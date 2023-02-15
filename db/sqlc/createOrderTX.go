package db

import (
	"context"
	"database/sql"
	"errors"

	ultils "github.com/bernie-pham/ecommercePlatform/ultil"
	"github.com/rs/zerolog/log"
)

type CreateOrderTXParams struct {
	UserID          int        `json:"user_id"`
	Items           []CartItem `json:"items"`
	DealID          int        `json:"deal_id"`
	AfterCreateFunc func(merchantID_OrderID_map map[int32]int64) error
}

func (store *SQLStore) CreateOrderTX_v2(ctx context.Context, arg CreateOrderTXParams) error {
	afterCreateArg := make(map[int32]int64)
	err := store.execTX(ctx, func(q *Queries) error {
		var base_total_price float32
		var discount_total_price float32
		// key is merchant_id
		// value is slice of order_item
		orderMap := make(map[int32][]CreateOrderItemParams)
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
			orderItemParam := CreateOrderItemParams{
				ProductEntryID: item.ProductEntryID,
				Quantity:       item.Quantity,
				TotalPrice:     total_price,
			}
			merchantID, err := q.GetMerchantIDbyPrID(ctx, product_entry.ProductID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to get deal")
				return err
			}
			orderMap[merchantID] = append(orderMap[merchantID], orderItemParam)
			// if _, ok := orderMap[merchantID]; ok {

			// } else {
			// 	orderMap[merchantID] = []CreateOrderItemParams{orderItemParam}
			// }
			err = q.DeleteCartItemByID(ctx, item.ID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to delete cart item")
				return err
			}
		}

		// Calculate price if order deal exists
		dealID, dealValid := ultils.GetOptionalInt(arg.DealID)
		if dealValid {
			deal, err := q.GetDealByID(ctx, int64(dealID))
			if err != nil && err != sql.ErrNoRows {
				log.Error().
					Err(err).
					Msg("failed to get deal")
				return err
			}
			orderDiscount := base_total_price * deal.DiscountRate
			if orderDiscount > float32(deal.DealLimit.Int32) {
				orderDiscount = float32(deal.DealLimit.Int32)
			}
			discount_total_price += orderDiscount
		}

		// after get calculate price and deal -> create order -> create merchant_order ->order item
		// Create order
		user_order_id, err := q.CreateOrderV2(ctx, CreateOrderV2Params{
			UserID: int64(arg.UserID),
			DealID: sql.NullInt64{
				Int64: int64(dealID),
				Valid: dealValid,
			},
			BasePrice:     base_total_price,
			DiscountPrice: discount_total_price,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create an order")
			return err
		}
		// create merchant order and order item
		for key, value := range orderMap {
			// at first create merchant_order
			arg := CreateMerchantOrderParams{
				MerchantID:  int64(key),
				TotalPrice:  1,
				OrderID:     user_order_id,
				OrderStatus: OrderStatusOpen,
			}
			merchant_order_id, err := q.CreateMerchantOrder(ctx, arg)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to create merchant order")
				return err
			}
			// next, create order items
			var totalPrice float32
			for _, item := range value {
				arg := CreateOrderItemV2Params{
					OrderID:         user_order_id,
					MerchantOrderID: merchant_order_id,
					Quantity:        item.Quantity,
					TotalPrice:      item.TotalPrice,
					ProductEntryID:  item.ProductEntryID,
				}
				err := q.CreateOrderItemV2(ctx, arg)
				if err != nil {
					log.Error().
						Err(err).
						Msg("failed to create order item")
					return err
				}
				totalPrice += item.TotalPrice
			}
			err = q.UpdateMerchantOrderTotalPrice(ctx, UpdateMerchantOrderTotalPriceParams{
				ID:         merchant_order_id,
				TotalPrice: totalPrice,
			})
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to update total price for merchant order")
				return err
			}
			afterCreateArg[key] = merchant_order_id
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = arg.AfterCreateFunc(afterCreateArg)
	return err
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
