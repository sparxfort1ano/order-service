package repository

import "time"

// Данные о доставке заказа
type Delivery struct {
	Name    string `json:"name"`    // имя получателя
	Phone   string `json:"phone"`   // телефон
	Zip     string `json:"zip"`     // индекс
	City    string `json:"city"`    // город
	Address string `json:"address"` // адрес
	Region  string `json:"region"`  // регион
	Email   string `json:"email"`   // email
}

// Данные об оплате заказа
type Payment struct {
	Transaction  string `json:"transaction"`   // номер транзакции
	RequestID    string `json:"request_id"`    // идентификатор запроса
	Currency     string `json:"currency"`      // валюта
	Provider     string `json:"provider"`      // провайдер оплаты
	Amount       int    `json:"amount"`        // общая сумма
	PaymentDt    int64  `json:"payment_dt"`    // дата оплаты
	Bank         string `json:"bank"`          // банк
	DeliveryCost int    `json:"delivery_cost"` // стоимость доставки
	GoodsTotal   int    `json:"goods_total"`   // стоимость товаров
	CustomFee    int    `json:"custom_fee"`    // таможенная пошлина
}

// Данные о конкретной товаре в заказе
type Item struct {
	ChrtID      int    `json:"chrt_id"`      // id товара
	TrackNumber string `json:"track_number"` // трек-номер
	Price       int    `json:"price"`        // цена
	Rid         string `json:"rid"`          // идентификатор позиции
	Name        string `json:"name"`         // название товара
	Sale        int    `json:"sale"`         // скидка
	Size        string `json:"size"`         // размер
	TotalPrice  int    `json:"total_price"`  // цена с учетом скидки
	NmID        int    `json:"nm_id"`        // артикул
	Brand       string `json:"brand"`        // бренд
	Status      int    `json:"status"`       // статус
}

// Полная структура заказа
type Order struct {
	OrderUID        string    `json:"order_uid"`          // уникальный id заказа
	TrackNumber     string    `json:"track_number"`       // трек-номер заказа
	Entry           string    `json:"entry"`              // точка входа
	Delivery        Delivery  `json:"delivery"`           // информация о доставке
	Payment         Payment   `json:"payment"`            // информация об оплате
	Items           []Item    `json:"items"`              // список товаров
	Locale          string    `json:"locale"`             // язык
	InternalSign    string    `json:"internal_signature"` // внутренняя подпись
	CustomerID      string    `json:"customer_id"`        // id покупателя
	DeliveryService string    `json:"delivery_service"`   // служба доставки
	ShardKey        string    `json:"shardkey"`           // ключ шардинга
	SmID            int       `json:"sm_id"`              // id склада
	DateCreated     time.Time `json:"date_created"`       // дата создания
	OofShard        string    `json:"oof_shard"`          // шард
}
