package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AddCartItemRequest struct {
	ProductEntryID int `json:"product_entry_id" binding:"required"`
	Quantity       int `json:"quantity" binding:"required"`
}

type CartItem struct {
	ItemID         int `json:"ItemID"`
	Quantity       int `json:"quantity"`
	ProductEntryID int `json:"product_entry_id"`
	ItemPrice
}

type ItemPrice struct {
	BasePrice int     `json:"base"`
	Discount  float64 `json:"discount"`
}

// AddCartItem add an item to cart and return new cart item list
func (server *Server) AddCartItem(ctx *gin.Context) {

	value, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		err := errors.New("invalid access token")
		log.Error().
			Err(err).
			Msg("failed to get userID from access payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	authPayload := value.(*token.Payload)

	var req AddCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// checking product_entry has enough quantity
	product_entry, err := server.store.GetPEntry(ctx, int64(req.ProductEntryID))
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to get product entry")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	if req.Quantity > int(product_entry.Quantity) {
		ctx.JSON(http.StatusNotAcceptable, errorResponse(ErrBadRequestParameter))
		return
	}

	// Get Base Price for product
	todayPrice, err := server.store.GetTodayBasePrice(ctx, product_entry.ProductID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to get product base price")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}

	// Get Discount if any
	var sale float32
	if product_entry.DealID.Valid {
		deal, err := server.store.GetDealByID(ctx, product_entry.DealID.Int64)

		if err != nil && err != sql.ErrNoRows {
			log.Error().
				Err(err).
				Msg("failed to get product deal price")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}

		sale = deal.DiscountRate * float32(todayPrice)
		if deal.DealLimit.Valid && sale > float32(deal.DealLimit.Int32) {
			sale = float32(deal.DealLimit.Int32)
		}
	}

	// Checking product entry id
	cart_id, err := server.store.GetCartItemByEntryID(ctx, db.GetCartItemByEntryIDParams{
		UserID:         int64(authPayload.UserID),
		ProductEntryID: product_entry.ID,
	})

	if err != nil && err != sql.ErrNoRows {
		log.Error().
			Err(err).
			Msg("failed to get cart item")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	var cartItem db.CartItem
	if cart_id != 0 {
		cartItem, err = server.store.UpdateCartItem(ctx, db.UpdateCartItemParams{
			ID:         cart_id,
			ModifiedAt: time.Now(),
			Quantity:   int32(req.Quantity),
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to update cart item")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		arg := db.AddCartItemParams{
			ProductEntryID: int64(req.ProductEntryID),
			Quantity:       int32(req.Quantity),
			UserID:         int64(authPayload.UserID),
		}

		cartItem, err = server.store.AddCartItem(ctx, arg)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to add item to cart")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}
	}

	itemPrice := ItemPrice{
		BasePrice: int(todayPrice),
		Discount:  float64(sale),
	}

	rsp := CartItem{
		ItemID:         int(cartItem.ID),
		ProductEntryID: int(cartItem.ProductEntryID),
		Quantity:       int(cartItem.Quantity),
		ItemPrice:      itemPrice,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// ListCartItems list all the items in user's cart
// url /carts
func (server *Server) ListCartItems(ctx *gin.Context) {
	// TODO: implement session management
	value, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		err := errors.New("invalid access token")
		log.Error().
			Err(err).
			Msg("failed to get userID from access payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	authPayload := value.(*token.Payload)

	var items []CartItem
	cartItems, err := server.store.ListCartItemsByUserID(ctx, int64(authPayload.UserID))
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to list item in cart")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
		return
	}
	fmt.Println("num item on cart:", authPayload.UserID)
	basePrices := make(map[int64]int32)
	dealMap := make(map[int64]db.Deal)
	for _, item := range cartItems {
		log.Info().Msgf("item %v", item)
		product_entry, err := server.store.GetPEntry(ctx, item.ProductEntryID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to get product entry")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}
		newQuantity := item.Quantity
		if newQuantity > product_entry.Quantity {
			newQuantity = 1
		}
		// Get base price
		var todayPrice int32
		if value, ok := basePrices[product_entry.ProductID]; ok {
			todayPrice = value
		} else {
			todayPrice, err = server.store.GetTodayBasePrice(ctx, product_entry.ProductID)
			basePrices[product_entry.ProductID] = todayPrice
		}
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to get product base price")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
			return
		}
		var sale float32
		if product_entry.DealID.Valid {
			var deal db.Deal
			if value, ok := dealMap[product_entry.DealID.Int64]; ok {
				deal = value
			} else {
				deal, err = server.store.GetDealByID(ctx, product_entry.DealID.Int64)
				dealMap[product_entry.DealID.Int64] = deal
			}
			if err != nil {
				if err != sql.ErrNoRows {
					log.Error().
						Err(err).
						Msg("failed to get product deal price")
					ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInteralErrServer))
					return
				}
			}

			sale = deal.DiscountRate * float32(todayPrice)
			if deal.DealLimit.Valid && sale > float32(deal.DealLimit.Int32) {
				sale = float32(deal.DealLimit.Int32)
			}
		}
		price := ItemPrice{
			BasePrice: int(todayPrice),
			Discount:  float64(sale),
		}
		items = append(items, CartItem{
			ItemID:         int(item.ID),
			Quantity:       int(newQuantity),
			ProductEntryID: int(product_entry.ID),
			ItemPrice:      price,
		})

	}
	ctx.JSON(http.StatusOK, items)
}
