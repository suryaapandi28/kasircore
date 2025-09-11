package entity

import "time"

type Tickets struct {
	Tickets_id     string    `json:"tickets_id"  gorm:"primarykey"`
	Transaction_id string    `json:"transaction_id"`
	Event_id       string    `json:"event_id"`
	Code_qr        string    `json:"code_qr"`
	Name_event     string    `json:"name_event"`
	Ticket_date    time.Time `json:"ticket_date"`
	Qty            string    `json:"qty"`
	Auditable
}

func NewTicket(tickets_id, transaction_id, event_id, code_qr, name_event, qty string) *Tickets {
	return &Tickets{
		Tickets_id:     tickets_id,
		Transaction_id: transaction_id,
		Event_id:       event_id,
		Code_qr:        code_qr,
		Name_event:     name_event,
		Ticket_date:    time.Now(),
		Qty:            qty,
		Auditable:      NewAuditable(),
	}
}
