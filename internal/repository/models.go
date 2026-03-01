package repository

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Order struct {
	OrderUid          string    `json:"order_uid" validate:"required,alphanum"`
	TrackNumber       string    `json:"track_number" validate:"required,alphanum"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,gt=0,dive"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required,alphanum"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int64     `json:"sm_id" validate:"required,number"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required,alphaspace"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip"`
	City    string `json:"city" validate:"required,alphaspace"`
	Address string `json:"address" validate:"required,alphanumspace"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"required,email"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" validate:"required,number"`
	TrackNumber string `json:"track_number" validate:"required,alphanum"`
	Price       int64  `json:"price" validate:"required,number"`
	Rid         string `json:"rid"`
	Name        string `json:"name" validate:"required,alphaspace"`
	Sale        int64  `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price" validate:"required,number"`
	NmID        int64  `json:"nm_id" validate:"required,number"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required,alphanum"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required,len=3"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount" validate:"required,number"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost" validate:"required,number"`
	GoodsTotal   int64  `json:"goods_total" validate:"required,number"`
	CustomFee    int64  `json:"custom_fee" validate:"required,number"`
}

func (o *Order) Validate() error {
	return validate.Struct(o)
}

type OrderRepository interface {
	// Data inserting into database
	SaveOrder(context.Context, *Order) error
	// Select data from database by orderUID
	GetOrderById(context.Context, string) (*Order, error)
}
