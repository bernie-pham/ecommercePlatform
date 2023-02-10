package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey   = "authorization"
	authorizationBearerType  = "bearer"
	authorizationPayloadKey  = "authorization_payload"
	payloadAccessTokenFooter = "access_token"
)

func authMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationBearerType {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, footer, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if footer != payloadAccessTokenFooter {
			err := errors.New("invalid access token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func merchantAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session_payload, ok := ctx.Get(authorizationPayloadKey)
		if !ok {
			err := errors.New("failed to get authorization payload in session")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		auth_payload := session_payload.(*token.Payload)
		if auth_payload.Access_level != 5 {
			err := errors.New("unauthorized user")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Next()
	}
}
