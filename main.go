package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/bernie-pham/ecommercePlatform/api"
	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mail "github.com/xhit/go-simple-mail/v2"
)

func main() {
	config, err := ultils.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load server config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// pure http server without any framework
	// listener, err := net.Listen("tcp", config.ServerAddr)
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("failed to listen on: %s", config.ServerAddr)
	// }

	// http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome to Ecommerce Platform"))
	// }))

	// Connect to database
	conn, err := sql.Open(config.DBDriverName, config.DBSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to DB")
	}

	// Init Session Token Maker
	token_maker, err := token.NewPasetoMaker(config.SymKey)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to init token maker")
	}
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddr,
	}
	// Init Asynq
	taskDistributor := async.NewRedisTaskDistributor(redisOpt)
	mailClient := newMailClient()
	// Create new db store
	store := db.NewStore(conn)
	go runTaskProccessor(redisOpt, store, mailClient, config.MailSender)
	runGinServer(config, store, token_maker, taskDistributor)

}

func runGinServer(
	config ultils.Config,
	store db.Store, maker token.TokenMaker,
	taskDistributor async.TaskDistributor,
) {
	server, err := api.NewServer(config, store, maker, taskDistributor)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot create Server")
	}
	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Cannot Listen Server on %s", config.HTTPServerAddr)
	}
}

func runTaskProccessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	mailClient *mail.SMTPClient,
	sender string,
) {
	taskProccessor := async.NewRedisTaskProccessor(redisOpt, store, mailClient, sender)
	log.Info().Msg("Start task processor")
	err := taskProccessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func newMailClient() *mail.SMTPClient {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot initialize mail client")
	}
	return client
}
