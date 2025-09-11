package entity

type Payments struct {
	Payment_id          string `json:"payment_id"  gorm:"primarykey"`
	Transaksi_id        string `json:"transaksi_id"`
	Status_pay          string `json:"status_pay"`
	Pay_time            string `json:"pay_time"`
	Pay_settlement_time string `json:"pay_settlement_time"`
	Pay_type            string `json:"pay_type"`
	Signature_key       string `json:"signature_key"`
	Auditable
}

func NewPaymentdata(payment_id, transaksi_id, status_pay, pay_time, pay_settlement_time, pay_type, signature_key string) *Payments {
	return &Payments{
		Payment_id:          payment_id,
		Transaksi_id:        transaksi_id,
		Status_pay:          status_pay,
		Pay_time:            pay_time,
		Pay_settlement_time: pay_settlement_time,
		Pay_type:            pay_type,
		Signature_key:       signature_key,
		Auditable:           NewAuditable(),
	}
}
