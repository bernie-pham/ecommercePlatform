package gapi

import (
	"fmt"

	"github.com/bernie-pham/ecommercePlatform/async"
	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/pb"
	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/bernie-pham/ecommercePlatform/token"
	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
)

type Server struct {
	pb.UnimplementedOrderManagementServer
	config          ultils.Config
	store           db.Store
	tokenMaker      token.TokenMaker
	taskDistributor async.TaskDistributor
	sessionRepo     session.SessionRepository
}

func NewServer(
	config ultils.Config,
	store db.Store,
	taskDistributor async.TaskDistributor,
	sessionRepo session.SessionRepository,
) (*Server, error) {
	// Return a pointer to Server instance because it avoids to clone this instance more than once.
	tokenMaker, err := token.NewPasetoMaker(config.SymKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
		sessionRepo:     sessionRepo,
	}
	return server, nil
}
