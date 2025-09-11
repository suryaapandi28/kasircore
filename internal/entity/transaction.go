package entity

import (
	"time"
)

type Transactions struct {
	Transactions_id string `json:"transactions_id"`
	Cart_id         string `json:"cart_id"`
	User_id         string `json:"user_id"`
	Fullname_user   string `json:"fullname_user"`
	Trx_date        time.Time
	Payment         string `json:"payment"`
	Payment_url     string `json:"payment_url"`
	Amount          string `json:"amount"`
	Status          string `json:"status"`
	Auditable
}

func NewTransaction(transactions_id, cart_id, user_id, fullname_user, payment, payment_url, amount, status string) *Transactions {
	return &Transactions{
		Transactions_id: transactions_id,
		Cart_id:         cart_id,
		User_id:         user_id,
		Fullname_user:   fullname_user,
		Trx_date:        time.Now(),
		Payment:         payment,
		Payment_url:     payment_url,
		Amount:          amount,
		Status:          status,
		Auditable:       NewAuditable(),
	}
}

func UpdateTransaction(transactions_id, status string) *Transactions {
	return &Transactions{
		Transactions_id: transactions_id,
		Status:          status,
		Auditable:       UpdateAuditable(),
	}
}

type Transaction_details struct {
	Transaction_details_id string `json:"transaction_details_id"`
	Transaction_id         string `json:"transaction_id"`
	Event_id               string `json:"event_id" `
	Name_event             string `json:"name_event"`
	Qty_event              string `json:"qty_event"`
	Price                  string `json:"price"`
	Ticket_date            string `json:"ticket_date"`
	Auditable
}

func NewTransactiondetail(transaction_details_id, event_id, transaction_id, name_event, qty_event, price, ticket_date string) *Transaction_details {
	return &Transaction_details{
		Transaction_details_id: transaction_details_id,
		Transaction_id:         transaction_id,
		Event_id:               event_id,
		Name_event:             name_event,
		Qty_event:              qty_event,
		Price:                  price,
		Ticket_date:            ticket_date,
		Auditable:              NewAuditable(),
	}
}
