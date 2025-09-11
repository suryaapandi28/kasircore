package entity

type Carts struct {
	Cart_id     string `json:"cart_id"  gorm:"primarykey"`
	User_id     string `json:"user_id"`
	Event_id    string `json:"event_id"`
	Qty         string `json:"qty"`
	Ticket_date string `json:"ticket_date"`
	Price       string `json:"price"`
	Auditable
}

func NewCart(cart_id, user_id, event_id, qty, ticket_date, price string) *Carts {
	return &Carts{
		Cart_id:     cart_id,
		User_id:     user_id,
		Event_id:    event_id,
		Qty:         qty,
		Ticket_date: ticket_date,
		Price:       price,
		Auditable:   NewAuditable(),
	}
}
