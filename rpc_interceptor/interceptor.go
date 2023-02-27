package interceptor

import (
	"context"

	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Interceptor struct {
	accessLevels map[string][]int
	tokenMaker   token.TokenMaker
	sessionRepo  session.SessionRepository
}

func NewInterceptor(accessLevels map[string][]int, tokenMaker token.TokenMaker, sessionRepo session.SessionRepository) *Interceptor {
	return &Interceptor{
		accessLevels: accessLevels,
		tokenMaker:   tokenMaker,
		sessionRepo:  sessionRepo,
	}
}

func (interceptor *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx, err := interceptor.authorizeWithRedis(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (interceptor *Interceptor) Stream(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Info().Str("Method", info.FullMethod).Msg("unary Interceptor")
	return handler(ctx, req)
}
