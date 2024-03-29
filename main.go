package main

import (
	"database/sql"
	"net"
	"os"
	"time"

	"github.com/bernie-pham/ecommercePlatform/api"
	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/gapi"
	"github.com/bernie-pham/ecommercePlatform/pb"
	interceptor "github.com/bernie-pham/ecommercePlatform/rpc_interceptor"
	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/bernie-pham/ecommercePlatform/token"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/elastic/go-elasticsearch/v8"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mail "github.com/xhit/go-simple-mail/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// Init Elastic Client
	// elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
	// 	Addresses: []string{config.ElasticAddr},
	// })
	elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticAddr},
	})
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to Elastic")
	}

	// Connect to database
	conn, err := sql.Open(config.DBDriverName, config.DBSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to DB")
	}
	// Init redis client
	redis := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   1,
	})
	sessionRepo := session.NewRedisSessionRepo(redis)

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

	// Create new db store
	store := db.NewStore(conn)

	// Init Asynq
	taskDistributor := async.NewRedisTaskDistributor(redisOpt, store)
	mailClient := newMailClient()

	go runTaskProccessor(redisOpt, store, mailClient, config.MailSender, *elasticClient)
	go runGinServer(config, store, token_maker, taskDistributor)
	runGRPCServer(config, store, taskDistributor, token_maker, sessionRepo)
}

func runGRPCServer(
	config ultils.Config,
	store db.Store,
	taskDistributor async.TaskDistributor,
	tokenMaker token.TokenMaker,
	sessionRepo session.SessionRepository,
) {
	server, err := gapi.NewServer(config, store, taskDistributor, sessionRepo)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot create grpc server")
	}
	// TODO: logging GRPC
	// Create new grpc server
	interceptor := interceptor.NewInterceptor(accessLevelEnum(), tokenMaker, sessionRepo)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
	pb.RegisterOrderManagementServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to create listener for grpc service")
	}
	log.Info().Msgf("start GRPC listener at: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start GRPC server")
	}
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
	elasticClient elasticsearch.Client,
) {
	taskProccessor := async.NewRedisTaskProccessor(redisOpt, store, mailClient, sender, elasticClient)
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

func accessLevelEnum() map[string][]int {
	service := "/pb.OrderManagement/"
	return map[string][]int{
		service + "UpdateOrderStatus": {1},
	}
}
