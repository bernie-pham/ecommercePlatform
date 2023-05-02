package async

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeSyncTag     = "elastic:sync_tag"
	TypeSyncProduct = "elastic:sync_product"
	SyncChunkSize   = 3
)

type SyncDataTask struct {
	IDs []string `json:"index_ids"`
}

func (distributor *RedisTaskDistributor) DistributeSyncAllTagDataTask(
	ctx context.Context,
	opt ...asynq.Option,
) error {

	// Get sequentially chunk of records from db
	start_cursor := 0
	isLast := false
	for !isLast {
		// TODO: In case of list tag_id failed, design a system to retry at point of failure.
		tagIDs, err := distributor.store.ListTagID(ctx, db.ListTagIDParams{
			Limit:  SyncChunkSize,
			Offset: int32(start_cursor),
		})
		if err != nil {
			return fmt.Errorf("failed to list tag ids: %w", err)
		}
		IDLists := SyncDataTask{
			IDs: tagIDs,
		}
		// TODO: In case of marshal payload failed, design a system to retry at point of failure.
		jPayload, err := json.Marshal(IDLists)
		if err != nil {
			return fmt.Errorf("failed to marshal task payload: %w", err)
		}
		asynqTask := asynq.NewTask(TypeSyncTag, jPayload, opt...)
		taskInfo, err := distributor.client.EnqueueContext(ctx, asynqTask, opt...)
		if err != nil {
			return fmt.Errorf("failed to enqueue task: %w", err)
		}
		log.Info().
			Str("type", asynqTask.Type()).
			Bytes("payload", asynqTask.Payload()).
			Str("queue", taskInfo.Queue).Int("max_retry", taskInfo.MaxRetry).Msg("enqueue task")

		// check if size of IDLists less than chunk size
		if start_cursor >= 1000 {
			isLast = true
		} else if len(tagIDs) < SyncChunkSize {
			isLast = true
		} else {
			start_cursor += SyncChunkSize
		}
	}

	return nil
}
func (distributor *RedisTaskDistributor) DistributeSyncAllProductDataTask(
	ctx context.Context,
	opt ...asynq.Option,
) error {

	// Get sequentially chunk of records from db
	start_cursor := 0
	isLast := false
	for !isLast {
		// TODO: In case of list tag_id failed, design a system to retry at point of failure.
		productIDs, err := distributor.store.ListProductID(ctx, db.ListProductIDParams{
			Limit:  SyncChunkSize,
			Offset: int32(start_cursor),
		})
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return fmt.Errorf("failed to list tag ids: %w", err)
		}
		IDLists := SyncDataTask{
			IDs: productIDs,
		}
		// TODO: In case of marshal payload failed, design a system to retry at point of failure.
		jPayload, err := json.Marshal(IDLists)
		if err != nil {
			return fmt.Errorf("failed to marshal task payload: %w", err)
		}
		asynqTask := asynq.NewTask(TypeSyncProduct, jPayload, opt...)
		taskInfo, err := distributor.client.EnqueueContext(ctx, asynqTask, opt...)
		if err != nil {
			return fmt.Errorf("failed to enqueue task: %w", err)
		}
		log.Info().
			Str("type", asynqTask.Type()).
			Bytes("payload", asynqTask.Payload()).
			Str("queue", taskInfo.Queue).Int("max_retry", taskInfo.MaxRetry).Msg("enqueue task")

		// check if size of IDLists less than chunk size
		if start_cursor >= 1000 {
			isLast = true
		} else if len(productIDs) < SyncChunkSize {
			isLast = true
		} else {
			start_cursor += SyncChunkSize
		}
	}

	return nil
}

type SyncNewPayload struct {
	IDs      []string `json:"index_ids"` // sync-all, sync-new
	TaskType string   `json:"task_type"`
}

// DistributeSyncNewDataTask distributes task for both new tag & new product data
func (distributor *RedisTaskDistributor) DistributeSyncNewDataTask(
	ctx context.Context,
	payload *SyncNewPayload,
	opt ...asynq.Option,
) error {

	// TODO: In case of marshal payload failed, design a system to retry at point of failure.
	isOutOfData := false
	start := 0
	var end int
	if SyncChunkSize <= len(payload.IDs) {
		end = SyncChunkSize
	} else {
		end = len(payload.IDs)
	}
	for !isOutOfData {
		newPayload := SyncDataTask{
			IDs: payload.IDs[start:end],
		}
		jPayload, err := json.Marshal(newPayload)
		if err != nil {
			return fmt.Errorf("failed to marshal task payload: %w", err)
		}
		asynqTask := asynq.NewTask(payload.TaskType, jPayload, opt...)
		taskInfo, err := distributor.client.EnqueueContext(ctx, asynqTask, opt...)
		if err != nil {
			return fmt.Errorf("failed to enqueue task: %w", err)
		}
		log.Info().
			Str("type", asynqTask.Type()).
			Bytes("payload", asynqTask.Payload()).
			Str("queue", taskInfo.Queue).Int("max_retry", taskInfo.MaxRetry).Msg("enqueue task")

		if len(payload.IDs) > end {
			start += SyncChunkSize
			if (end + SyncChunkSize) >= len(payload.IDs) {
				end = len(payload.IDs)
			} else {
				end += SyncChunkSize
			}
		} else {
			isOutOfData = true
		}

	}

	return nil
}
