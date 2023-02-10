package async

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	TypeEmailDeliver = "email:deliver"
	TypeNotification = "notification:notify"
)

type EmailDeliveryPayload struct {
	Email        string `json:"email_recipient"`
	Url          string `json:"email_url"`
	Subject      string `json:"email_subject"`
	Msg          string `json:"email_msg"`
	EmailTemplte string `json:"email_tmpl"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendMail(
	ctx context.Context,
	payload *EmailDeliveryPayload,
	opt ...asynq.Option,
) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TypeEmailDeliver, jsonPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueue task")
	return nil
}

func (handler *RedisTaskProccessor) HandleEmailDeliveryTask(ctx context.Context, task *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error().
			Err(err).
			Msg("failed to handle Email Delivery task")
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	err := sendMail(handler.mailSender, payload, handler.mailClient)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to handle Email Delivery task")
		return fmt.Errorf("sendMail failed: %v", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", payload.Email).Msg("processed send email task")
	return nil
}

func sendMail(sender string, payload EmailDeliveryPayload, mailClient *mail.SMTPClient) error {
	email := mail.NewMSG()
	email.SetFrom(sender).
		AddTo(payload.Email).
		SetSubject(payload.Subject)
	mailTmpl, err := ioutil.ReadFile(payload.EmailTemplte)
	if err != nil {
		return fmt.Errorf("failed to load emal template: %v", err)
	}
	mailString := string(mailTmpl)
	mailSubject := strings.Replace(mailString, "[%subject%]", payload.Subject, 1)
	mailBody := strings.Replace(mailSubject, "[%body%]", payload.Msg, 1)
	mailURL := strings.Replace(mailBody, "[%url%]", payload.Url, 1)
	email.SetBody(mail.TextHTML, mailURL)
	err = email.Send(mailClient)
	if err != nil {
		return errors.New("failed to send email")
	}
	return nil
}

type NotificationPayload struct {
	RecipientID int    `json:"recipient_id"`
	Msg         string `json:"msg"`
	Title       string `json:"title"`
}

func (distributor *RedisTaskDistributor) DistributeTaskNotification(
	ctx context.Context,
	payload *NotificationPayload,
	opt ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TypeNotification, jsonPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task, opt...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueue task")
	return nil
}

func (processor *RedisTaskProccessor) HandleTaskNotification(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload NotificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	err := notify(ctx, &payload, processor.store)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to handle Notification task")
		return err
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Int("notify", payload.RecipientID).Msg("processed notification task")
	return nil
}

func notify(ctx context.Context, payload *NotificationPayload, store db.Store) error {
	arg := db.CreateNotificationParams{
		Message:     payload.Msg,
		RecipientID: int64(payload.RecipientID),
		Title:       payload.Title,
	}
	_, err := store.CreateNotification(ctx, arg)
	if err != nil {
		return fmt.Errorf("notify failed: %v", err)
	}
	return nil
}
