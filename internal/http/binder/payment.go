package binder

type Payments struct {
	Payment_id          string `json:"payment_id"  validate:"required,payment_id"`
	Transaksi_id        string `json:"transaksi_id" validate:"required,transaksi_id"`
	Status_pay          string `json:"status_pay" validate:"required,status_pay"`
	Pay_time            string `json:"pay_time" validate:"required,pay_time"`
	Pay_settlement_time string `json:"pay_settlement_time" validate:"required,pay_settlement_time"`
	Pay_type            string `json:"pay_type" validate:"required,pay_type"`
	Signature_key       string `json:"signature_key" validate:"required,signature_key"`
}
