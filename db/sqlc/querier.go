// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddCartItem(ctx context.Context, arg AddCartItemParams) (CartItem, error)
	CreateDeal(ctx context.Context, arg CreateDealParams) (Deal, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)
	CreateMerchantOrder(ctx context.Context, arg CreateMerchantOrderParams) (int64, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (int64, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreateOrderItemV2(ctx context.Context, arg CreateOrderItemV2Params) error
	CreateOrderV2(ctx context.Context, arg CreateOrderV2Params) (int64, error)
	CreatePColour(ctx context.Context, colourName string) (ProductColour, error)
	CreatePCriteria(ctx context.Context, criteria string) (ProductGeneralCriterium, error)
	CreatePEntry(ctx context.Context, arg CreatePEntryParams) (ProductEntry, error)
	CreatePPrice(ctx context.Context, arg CreatePPriceParams) (ProductPricing, error)
	CreatePSize(ctx context.Context, sizeValue string) (ProductSize, error)
	CreateProTag(ctx context.Context, arg CreateProTagParams) (ProductTagsProduct, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateTag(ctx context.Context, name string) (ProductTag, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVerification(ctx context.Context, arg CreateVerificationParams) (Verification, error)
	DeleteAllCartItemByUserID(ctx context.Context, userID int64) error
	DeleteCartItemByID(ctx context.Context, id int64) error
	DeleteOrderItem(ctx context.Context, arg DeleteOrderItemParams) error
	DeletePColour(ctx context.Context, id int64) error
	DeletePCriteria(ctx context.Context, id int64) error
	DeletePSize(ctx context.Context, id int64) error
	DeleteProTag(ctx context.Context, arg DeleteProTagParams) error
	DeleteTag(ctx context.Context, id int64) error
	DisableMerchant(ctx context.Context, id int64) error
	EnableMerchant(ctx context.Context, id int64) error
	GetCartItemByEntryID(ctx context.Context, arg GetCartItemByEntryIDParams) (int64, error)
	GetCartItemByID(ctx context.Context, id int64) (CartItem, error)
	GetCurrentOrder(ctx context.Context, userID int64) (Order, error)
	GetDealByID(ctx context.Context, id int64) (Deal, error)
	GetMerchant(ctx context.Context, id int64) (Merchant, error)
	GetMerchantByCartID(ctx context.Context, id int64) (sql.NullInt32, error)
	GetMerchantIDByPEntry(ctx context.Context, id int64) (sql.NullInt32, error)
	GetMerchantIDbyPrID(ctx context.Context, id int64) (int32, error)
	GetMerchantOrder(ctx context.Context, arg GetMerchantOrderParams) (MerchantOrder, error)
	GetPEntry(ctx context.Context, id int64) (ProductEntry, error)
	GetSession(ctx context.Context, id uuid.UUID) (GetSessionRow, error)
	GetTodayBasePrice(ctx context.Context, productID int64) (int32, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetVerification(ctx context.Context, id uuid.UUID) (Verification, error)
	ListActivePEntriesByPID(ctx context.Context, productID int64) ([]ProductEntry, error)
	ListAllProducts(ctx context.Context) ([]Product, error)
	ListCartItemsByUserID(ctx context.Context, userID int64) ([]CartItem, error)
	ListDealsByMerchantID(ctx context.Context, merchantID int64) ([]Deal, error)
	ListMerchants(ctx context.Context, userID int64) ([]Merchant, error)
	ListNotifications(ctx context.Context, recipientID int64) ([]Notification, error)
	ListOrder(ctx context.Context, userID int64) ([]Order, error)
	ListOrderItemsByMerchantOrderID(ctx context.Context, merchantOrderID int64) ([]ListOrderItemsByMerchantOrderIDRow, error)
	ListPColours(ctx context.Context) ([]ProductColour, error)
	ListPCriterias(ctx context.Context) ([]ProductGeneralCriterium, error)
	ListPEntriesByPID(ctx context.Context, productID int64) ([]ProductEntry, error)
	ListPSizes(ctx context.Context) ([]ProductSize, error)
	ListPriceByPID(ctx context.Context, productID int64) ([]ProductPricing, error)
	ListProTags(ctx context.Context) ([]ProductTagsProduct, error)
	ListProductByMerchantID(ctx context.Context, merchantID int32) ([]Product, error)
	ListProductTags(ctx context.Context, productsID int64) ([]ProductTag, error)
	ListProductsByTagID(ctx context.Context, productTagsID int64) ([]Product, error)
	ListTags(ctx context.Context) ([]ProductTag, error)
	OccupyVerification(ctx context.Context, id uuid.UUID) error
	UpdateCartItem(ctx context.Context, arg UpdateCartItemParams) (CartItem, error)
	UpdateEntryQuantity(ctx context.Context, arg UpdateEntryQuantityParams) error
	UpdateMerchant(ctx context.Context, arg UpdateMerchantParams) (Merchant, error)
	UpdateMerchantOrderStatus(ctx context.Context, arg UpdateMerchantOrderStatusParams) error
	UpdateMerchantOrderTotalPrice(ctx context.Context, arg UpdateMerchantOrderTotalPriceParams) error
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateOrderItemQuantity(ctx context.Context, arg UpdateOrderItemQuantityParams) (OrderItem, error)
	UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error
	UpdatePEntry(ctx context.Context, arg UpdatePEntryParams) (ProductEntry, error)
	UpdatePPrice(ctx context.Context, arg UpdatePPriceParams) (ProductPricing, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
