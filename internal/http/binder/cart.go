package binder

type AddCartRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	EventId string `json:"event_id" validate:"required"`
	Qty     int    `json:"qty" validate:"required"`
}

type UpdateQuantityLessRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	EventId string `json:"event_id" validate:"required"`
}
type UpdateQuantityAddRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	EventId string `json:"event_id" validate:"required"`
}

type GetCartResponse struct {
	EventID      string `json:"event_id"`
	TitleEvent   string `json:"title_event"`
	CityEvent    string `json:"city_event"`
	AddressEvent string `json:"address_event"`
	Qty          int    `json:"qty"`
	TicketDate   string `json:"ticket_date"`
	Price        int    `json:"price"`
}

type FindCartByIdRequest struct {
	CartID string `param:"id" validate:"required"`
}

type FindCartByUserIdRequest struct {
	UserID string `param:"id" validate:"required"`
}

type RemoveCartRequest struct {
	CartID string `param:"id" validate:"required"`
}
