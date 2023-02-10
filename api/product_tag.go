package api

import (
	"fmt"
	"net/http"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateTagRequest struct {
	Tag string `json:"tag_name"`
}

func (server *Server) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	Tag, err := server.store.CreateTag(ctx, req.Tag)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tag")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, Tag)
}

func (server *Server) ListTags(ctx *gin.Context) {
	Tags, err := server.store.ListTags(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list tags")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, Tags)
}

type CreateProtagRequest struct {
	ProductID int `json:"product_id"`
	TagID     int `json:"tag_id"`
}

func (server *Server) CreateProtag(ctx *gin.Context) {
	var req CreateProtagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateProTagParams{
		ProductTagsID: int64(req.ProductID),
		ProductsID:    int64(req.ProductID),
	}

	Protag, err := server.store.CreateProTag(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create product colour")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, Protag)
}

func (server *Server) ListProtags(ctx *gin.Context) {
	Protags, err := server.store.ListProTags(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create product colour")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, Protags)
}

type ListProductWithTagRequest struct {
	TagID int `uri:"id" binding:"required,min=1"`
}

func (server *Server) ListProductsWithTag(ctx *gin.Context) {
	var req ListProductWithTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	fmt.Println("Tag ID ", req.TagID)
	products, err := server.store.ListProductsByTagID(ctx, int64(req.TagID))
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to list products having tag_id")
	}
	ctx.JSON(http.StatusOK, products)
}
