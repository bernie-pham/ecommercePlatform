package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateProductEntryParam struct {
	ProductID string `json:"pro_id" binding:"required"`
	ColourID  int    `json:"colour_id"`
	SizeID    int    `json:"size_id"`
	GeneralID int    `json:"general_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	DealID    int    `json:"deal_id"`
}

func (server *Server) CreateProductEntry(ctx *gin.Context) {
	var req CreateProductEntryParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	colourID, isColorID := getOptionalInt(req.ColourID)
	sizeID, isSizeID := getOptionalInt(req.SizeID)
	genID, isGenID := getOptionalInt(req.GeneralID)
	dealID, isDealID := getOptionalInt(req.DealID)

	arg := db.CreatePEntryParams{
		ProductID: req.ProductID,
		ColourID: sql.NullInt64{
			Int64: int64(colourID),
			Valid: isColorID,
		},
		SizeID: sql.NullInt64{
			Int64: int64(sizeID),
			Valid: isSizeID,
		},
		GeneralCriteriaID: sql.NullInt64{
			Int64: int64(genID),
			Valid: isGenID,
		},
		DealID: sql.NullInt64{
			Int64: int64(dealID),
			Valid: isDealID,
		},
		Quantity: int32(req.Quantity),
	}
	pentry, err := server.store.CreatePEntry(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create product entry")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, pentry)
}

type UpdateProductEntryParam struct {
	ID       int  `json:"id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
	DealID   int  `json:"deal_id"`
	IsActive bool `json:"is_active"`
}

func (server *Server) UpdateProductEntry(ctx *gin.Context) {
	var req UpdateProductEntryParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequestParameter))
		return
	}
	dealID, isDealID := getOptionalInt(req.DealID)
	quantity, isQuantity := getOptionalInt(req.Quantity)
	isModified := false
	isModified = isDealID || isModified
	isModified = isQuantity || isModified
	arg := db.UpdatePEntryParams{
		ID: int64(req.ID),
		Quantity: sql.NullInt32{
			Int32: int32(quantity),
			Valid: isQuantity,
		},
		DealID: sql.NullInt64{
			Int64: int64(dealID),
			Valid: isDealID,
		},
		IsActive: req.IsActive,
		ModifiedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: isModified,
		},
	}
	pentry, err := server.store.UpdatePEntry(ctx, arg)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to update product entry")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, pentry)
}

type ListProductEntriesReq struct {
	ID string `uri:"id" binding:"required"`
}

// ListProductEntries take uri parameter GET /product/:id
// return a list of entries
func (server *Server) ListProductEntries(ctx *gin.Context) {
	var req ListProductEntriesReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	entries, err := server.store.ListPEntriesByPID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		log.Error().
			Err(err).
			Msg("failed to get product entries")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	ctx.JSON(http.StatusOK, entries)
}
