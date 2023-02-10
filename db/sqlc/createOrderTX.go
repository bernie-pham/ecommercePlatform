package db

import (
	"context"
	"database/sql"
	"errors"

	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/rs/zerolog/log"
)

type CreateOrderTXParams struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
	DealID int        `json:"deal_id"`
}

func (store *SQLStore) CreateOrderTX(ctx context.Context, arg CreateOrderTXParams) error {
	err := store.execTX(ctx, func(q *Queries) error {
		orderID, err := q.CreateOrder(ctx, CreateOrderParams{
			UserID:        int64(arg.UserID),
			BasePrice:     0,
			DiscountPrice: 0,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create an order")
			return err
		}
		var base_total_price float32
		var discount_total_price float32
		for _, item := range arg.Items {

			product_entry, err := q.GetPEntry(ctx, item.ProductEntryID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to get product entry")
				return err
			}
			// checking quantity
			if product_entry.Quantity < item.Quantity {
				err = errors.New("invalid order item's quantity")
				log.Error().
					Err(err).
					Msg("product entry's quantity is not sufficient for order's quantity")
				return err
			}
			// checking for price
			price, err := q.GetTodayBasePrice(ctx, product_entry.ProductID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to get product base price")
				return err
			}

			var total_price float32
			total_price = float32(price * item.Quantity)
			base_total_price += total_price
			// checking any sale on this product entry
			if product_entry.DealID.Valid {
				deal, err := q.GetDealByID(ctx, product_entry.DealID.Int64)
				if err != nil && err != sql.ErrNoRows {
					log.Error().
						Err(err).
						Msg("failed to get deal")
					return err
				}
				discount := float32(price) * deal.DiscountRate
				// get total discount on single product entry
				if discount > float32(deal.DealLimit.Int32) {
					discount = float32(deal.DealLimit.Int32)
				}
				total_discount := discount * float32(item.Quantity)
				discount_total_price += total_discount
				total_price = total_price - total_discount
			}

			_, err = q.CreateOrderItem(ctx, CreateOrderItemParams{
				OrderID:        orderID,
				ProductEntryID: item.ProductEntryID,
				Quantity:       item.Quantity,
				TotalPrice:     total_price,
			})
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to create order item")
				return err
			}
			err = q.DeleteCartItemByID(ctx, item.ID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to delete cart item")
				return err
			}
		}

		dealID, valid := ultils.GetOptionalInt(arg.DealID)
		var orderDiscount float32
		if valid {
			deal, err := q.GetDealByID(ctx, int64(dealID))
			if err != nil && err != sql.ErrNoRows {
				log.Error().
					Err(err).
					Msg("failed to get deal")
				return err
			}
			orderDiscount = base_total_price * deal.DiscountRate
			if orderDiscount > float32(deal.DealLimit.Int32) {
				orderDiscount = float32(deal.DealLimit.Int32)
			}
			discount_total_price += orderDiscount
		}

		_, err = q.UpdateOrder(ctx, UpdateOrderParams{
			Status: NullOrderStatus{
				OrderStatus: OrderStatusOpen,
				Valid:       true,
			},
			BasePrice: sql.NullFloat64{
				Float64: float64(base_total_price),
				Valid:   true,
			},
			DiscountPrice: sql.NullFloat64{
				Float64: float64(discount_total_price),
				Valid:   true,
			},
			DealID: sql.NullInt64{
				Int64: int64(dealID),
				Valid: valid,
			},
			ID: orderID,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to update order")
			return err
		}
		return nil
	})
	return err
}
