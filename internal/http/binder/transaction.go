package binder

type TrasactionCreateRequest struct {
	Transactions_id string `json:"transactions_id" validate:"required,transactions_id"`
	Cart_id         string `json:"cart_id" validate:"required,cart_id"`
	User_id         string `json:"user_id" validate:"required,user_id"`
	Fullname_user   string `json:"fullname_user" validate:"required,fullname_user"`
	Trx_date        string `json:"trx_date" validate:"required,trx_date"`
	Payment         string `json:"payment" validate:"required,payment"`
	Payment_url     string `json:"payment_url" validate:"required,payment_url"`
	Amount          string `json:"amount" validate:"required,amount"`
	Status          string `json:"status" validate:"required,status"`
}

type TrasactionCreatedetailRequest struct {
	Transaction_details_id string `json:"transaction_details_id" validate:"required,transaction_details_id"`
	Transaction_id         string `json:"transaction_id" validate:"required,transaction_id"`
	Event_id               string `json:"event_id" validate:"required,event_id"`
	Name_event             string `json:"name_event" validate:"required,name_event"`
	Qty_event              string `json:"qty_event" validate:"required,qty_event"`
	Price                  string `json:"price" validate:"required,price"`
	Ticket_date            string `json:"ticket_date" validate:"required,ticket_date"`
}

type EventFindByIDRequest struct {
	Event_id string `param:"event_id" validate:"required"`
}

type CheckTrxFindByIDRequest struct {
	Transactions_id string `param:"transactions_id" validate:"required"`
}

type GetAllRequest struct {
	Key             string `json:"key" validate:"required,key"`
	Transactions_id string `json:"transactions_id" validate:"required,transactions_id"`
	User_id         string `json:"user_id" validate:"required,user_id"`
}
