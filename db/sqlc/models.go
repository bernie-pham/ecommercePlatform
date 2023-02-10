// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusArchived OrderStatus = "archived"
	OrderStatusCanceled OrderStatus = "canceled"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus
	Valid       bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.OrderStatus, nil
}

type ProductStatus string

const (
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
	ProductStatusInStock    ProductStatus = "in_stock"
	ProductStatusRunningLow ProductStatus = "running_low"
)

func (e *ProductStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProductStatus(s)
	case string:
		*e = ProductStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ProductStatus: %T", src)
	}
	return nil
}

type NullProductStatus struct {
	ProductStatus ProductStatus
	Valid         bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProductStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ProductStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProductStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProductStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.ProductStatus, nil
}

type CartItem struct {
	ID             int64     `json:"id"`
	ProductEntryID int64     `json:"product_entry_id"`
	Quantity       int32     `json:"quantity"`
	CreatedAt      time.Time `json:"created_at"`
	UserID         int64     `json:"user_id"`
	ModifiedAt     time.Time `json:"modified_at"`
}

type Country struct {
	Code          int32  `json:"code"`
	Name          string `json:"name"`
	ContinentName string `json:"continent_name"`
}

type Deal struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Code         sql.NullString `json:"code"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      time.Time      `json:"end_date"`
	Type         string         `json:"type"`
	DiscountRate float32        `json:"discount_rate"`
	MerchantID   int64          `json:"merchant_id"`
	DealLimit    sql.NullInt32  `json:"deal_limit"`
}

type Merchant struct {
	ID           int64     `json:"id"`
	CountryCode  int32     `json:"country_code"`
	MerchantName string    `json:"merchant_name"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       int64     `json:"user_id"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
}

type Notification struct {
	ID          int64     `json:"id"`
	Message     string    `json:"message"`
	RecipientID int64     `json:"recipient_id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
}

type Order struct {
	ID            int64          `json:"id"`
	UserID        int64          `json:"user_id"`
	Status        OrderStatus    `json:"status"`
	CreatedAt     sql.NullString `json:"created_at"`
	DealID        sql.NullInt64  `json:"deal_id"`
	BasePrice     float32        `json:"base_price"`
	DiscountPrice float32        `json:"discount_price"`
}

type OrderItem struct {
	OrderID        int64   `json:"order_id"`
	ProductEntryID int64   `json:"product_entry_id"`
	Quantity       int32   `json:"quantity"`
	TotalPrice     float32 `json:"total_price"`
}

type Product struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	MerchantID int32             `json:"merchant_id"`
	Status     NullProductStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

type ProductColour struct {
	ID         int64  `json:"id"`
	ColourName string `json:"colour_name"`
}

type ProductEntry struct {
	ID                int64         `json:"id"`
	ProductID         int64         `json:"product_id"`
	ColourID          sql.NullInt64 `json:"colour_id"`
	SizeID            sql.NullInt64 `json:"size_id"`
	GeneralCriteriaID sql.NullInt64 `json:"general_criteria_id"`
	Quantity          int32         `json:"quantity"`
	DealID            sql.NullInt64 `json:"deal_id"`
	IsActive          bool          `json:"is_active"`
	ModifiedAt        time.Time     `json:"modified_at"`
	CreatedAt         time.Time     `json:"created_at"`
}

type ProductGeneralCriterium struct {
	ID       int64  `json:"id"`
	Criteria string `json:"criteria"`
}

type ProductPricing struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	BasePrice int32     `json:"base_price"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	Priority  int32     `json:"priority"`
}

type ProductSize struct {
	ID        int64  `json:"id"`
	SizeValue string `json:"size_value"`
}

type ProductTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductTagsProduct struct {
	ProductTagsID int64 `json:"product_tags_id"`
	ProductsID    int64 `json:"products_id"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	HashedPassword    string    `json:"hashed_password"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
	AccessLevel       int32     `json:"access_level"`
}

type Verification struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	IsOccurpied bool      `json:"is_occurpied"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}