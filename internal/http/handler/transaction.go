package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
	"github.com/suryaapandi28/kasircore/pkg/token"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	paymentService     service.PaymentService
	tokenUseCase       token.TokenUseCase
	// paymentService     service.PaymentService
}

type CustomValidator struct {
	Validator *validator.Validate
}

type Responsemeta struct {
	Message string
	Status  bool
}

type TrasactionCreateRequestdata struct {
	Cart_id       string `json:"cart_id" validate:"required,cart_id"`
	User_id       string `json:"user_id" validate:"required,user_id"`
	Fullname_user string `json:"fullname_user" validate:"required,fullname_user"`
	Trx_date      string `json:"trx_date" validate:"required,trx_date"`
	Payment       string `json:"payment" validate:"required,payment"`
	Payment_url   string `json:"payment_url" validate:"required,payment_url"`
	Amount        string `json:"amount" validate:"required,amount"`
	Status        string `json:"status" validate:"required,status"`
}

type Paymentsdata struct {
	Order_id           string `json:"order_id"`
	Transaksi_id       string `json:"transaksi_id"`
	Transaction_status string `json:"transaction_status"`
	Transaction_time   string `json:"transaction_time"`
	Settlement_time    string `json:"settlement_time"`
	Payment_type       string `json:"payment_type"`
	Signature_key      string `json:"signature_key"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
	// Optionally, you could return the error to give each route more control over the status code
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (h *TransactionHandler) claimsjwtdata(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusBadRequest, "you must login first"))
	}
	claimsjwt := user.Claims.(*token.JwtCustomClaims)
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Create Data New Transaction", claimsjwt))
}

func NewTransactionHandler(transactionService service.TransactionService, tokenUseCase token.TokenUseCase, paymentService service.PaymentService) TransactionHandler {
	return TransactionHandler{transactionService: transactionService, tokenUseCase: tokenUseCase, paymentService: paymentService}

}

func calculatePPN(amount float64, tarifPPN float64) float64 {
	return amount * (tarifPPN / 100.0)
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	input := binder.TrasactionCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.
			JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input yang dimasukkan salah"))
	}
	if input.Cart_id == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
	}
	if input.User_id == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
	}
	if isValidUUID(input.Cart_id) {
		if isValidUUID(input.User_id) {

			Cart_id := uuid.MustParse(input.Cart_id)

			cartdata, err := h.transactionService.FindCartByID(Cart_id)

			if err != nil {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			}
			if cartdata.Cart_id == "" {
				return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found no Cart data"))
			}

			Event_id := uuid.MustParse(cartdata.Event_id)

			eventdata, err := h.transactionService.FindEventByID(Event_id)

			if err != nil {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			}

			if uuid.UUID(eventdata.Event_id) == uuid.Nil {
				return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found Event no data"))
			}

			User_id := uuid.MustParse(input.User_id)

			userdata, err := h.transactionService.FindUserByID(User_id)

			if err != nil {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			}

			if userdata.User_id == "" {
				return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found no data"))
			}
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusBadRequest, "you must login first"))
			}
			claimsjwt := user.Claims.(*token.JwtCustomClaims)
			if claimsjwt.Role == "user" {

				qtytrx := cartdata.Qty
				pricetrx := eventdata.Price_event

				amount2, err := strconv.Atoi(qtytrx)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! We failed to convert"))
				}

				amount := float64(pricetrx * amount2)
				ppn := calculatePPN(amount, 11)
				amounttotal := int((amount + ppn))

				amountfinal := strconv.Itoa(amounttotal)

				trx_id := uuid.New().String()
				status := "pending"

				url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
				// "enabled_payments": ["bca_va"],
				data := map[string]interface{}{
					"transaction_details": map[string]interface{}{
						"order_id":     trx_id,
						"gross_amount": int64(amounttotal),
					},

					"enabled_payments": []string{input.Payment},
					"customer_details": map[string]interface{}{
						"first_name": userdata.Fullname,
						"last_name":  "",
						"email":      userdata.Email,
						"phone":      "0" + userdata.Phone,
						"billing_address": map[string]interface{}{
							"first_name":   userdata.Fullname,
							"last_name":    "",
							"email":        userdata.Email,
							"phone":        strconv.Itoa(0) + userdata.Phone,
							"address":      "Indanesia",
							"city":         "Jakarta",
							"postal_code":  "12190",
							"country_code": "IDN",
						},
					},
				}

				// Mengubah data menjadi format JSON
				payload, err := json.Marshal(data)
				if err != nil {
					log.Fatalf("Failed to marshal JSON: %v", err)
				}

				// Membuat request HTTP POST
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
				if err != nil {
					log.Fatalf("Failed to create HTTP request: %v", err)
				}

				// Menambahkan header Content-Type: application/json
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Basic U0ItTWlkLXNlcnZlci1kazQ1S0Zpb21QRW9UajFqeWpiWWd1Z1k6Og==")

				// Membuat klien HTTP
				client := &http.Client{}

				// Mengirim request ke server
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("Failed to send HTTP request: %v", err)
				}
				defer resp.Body.Close()

				// Membaca respons dari server
				var snapRequest map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&snapRequest)
				snapRequesturl := snapRequest["redirect_url"]
				paymenturl := snapRequesturl.(string)
				if err != nil {
					log.Fatalf("Failed to decode JSON response: %v", err)
				}

				NewTrx := entity.NewTransaction(trx_id, input.Cart_id, input.User_id, userdata.Fullname, input.Payment, paymenturl, amountfinal, status)

				if NewTrx.Cart_id == "" {
					return c.JSON(http.StatusUnprocessableEntity, response.Errorfieldempty(http.StatusUnprocessableEntity, "Card_id"))
				}
				if NewTrx.User_id == "" {
					return c.JSON(http.StatusUnprocessableEntity, response.Errorfieldempty(http.StatusUnprocessableEntity, "User_id"))
				}
				if NewTrx.Payment == "" {
					return c.JSON(http.StatusUnprocessableEntity, response.Errorfieldempty(http.StatusUnprocessableEntity, "Payment"))
				}

				qtyevent := cartdata.Qty
				priceevent := eventdata.Price_event

				event2, err := strconv.Atoi(qtyevent)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! We failed to convert"))
				}

				price := priceevent * event2
				pricefinal := strconv.Itoa(price)

				trx_iddetail := uuid.New().String()

				NewTrxdetail := entity.NewTransactiondetail(trx_iddetail, cartdata.Event_id, trx_id, eventdata.Title_event, cartdata.Qty, pricefinal, cartdata.Ticket_date)

				transaction, err := h.transactionService.CreateTransaction(NewTrx)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}
				transactiondetail, err := h.transactionService.CreateTransactiondetail(NewTrxdetail)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}

				Pembayaran := make(map[string]interface{})

				Pembayaran["Payment_URL"] = snapRequesturl
				Pembayaran["Payment_Bank"] = input.Payment
				Pembayaran["Gross_Amount"] = amounttotal
				Pembayaran["PPN"] = ppn
				Pembayaran["Amount"] = amount

				trxpay := make(map[string]interface{})

				trxpay["Transaction_ID"] = transaction.Transactions_id
				trxpay["Transaction_Detail_ID"] = transactiondetail.Transaction_details_id
				trxpay["Pembayaran"] = Pembayaran

				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Create Data New Transaction", trxpay))

			} else if claimsjwt.Role == "admin" {
				return c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "Sorry! access denied"))
			} else {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! we cannot recognize your account"))
			}
		} else {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
		}
	} else {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
	}
}

func (h *TransactionHandler) CheckPayTransaction(c echo.Context) error {
	var input binder.CheckTrxFindByIDRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if input.Transactions_id == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
	}

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusBadRequest, "you must login first"))
	}
	claimsjwt := user.Claims.(*token.JwtCustomClaims)
	if claimsjwt.Role == "user" {

		if isValidUUID(input.Transactions_id) {

			transactions_id := uuid.MustParse(input.Transactions_id)

			transaction, err := h.transactionService.FindTrxByID(transactions_id)

			if err != nil {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			}
			if transaction.Transactions_id == "" {
				return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found no data"))
			}
			// transactions_id_checkpay := uuid.MustParse(input.Transactions_id)
			transactions_id_checkpay := transactions_id.String()

			url := "https://api.sandbox.midtrans.com/v2/" + transactions_id_checkpay + "/status"
			// "enabled_payments": ["bca_va"],
			data := map[string]interface{}{}

			// Mengubah data menjadi format JSON
			payload, err := json.Marshal(data)
			if err != nil {
				log.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Membuat request HTTP POST
			req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
			if err != nil {
				log.Fatalf("Failed to create HTTP request: %v", err)
			}

			// Menambahkan header Content-Type: application/json
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Basic U0ItTWlkLXNlcnZlci1kazQ1S0Zpb21QRW9UajFqeWpiWWd1Z1k6Og==")

			// Membuat klien HTTP
			client := &http.Client{}

			// Mengirim request ke server
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Failed to send HTTP request: %v", err)
			}
			defer resp.Body.Close()

			// Membaca respons dari server
			var checkpayreq map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&checkpayreq)
			sttscode := checkpayreq["status_code"]
			trxpay := checkpayreq["transaction_status"]

			if err != nil {
				log.Fatalf("Failed to decode JSON response: %v", err)
			}

			if sttscode == "404" {
				payreload := make(map[string]interface{})

				payreload["Payment_URL"] = transaction.Payment_url
				payreload["Payment_Bank"] = transaction.Payment
				payreload["Message"] = "Payment In Process"

				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Check Pay", payreload))
			} else if trxpay == "pending" {
				payreload := make(map[string]interface{})

				payreload["Payment_URL"] = transaction.Payment_url
				payreload["Payment_Bank"] = transaction.Payment
				payreload["Message"] = "Payment Pending"

				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Check Pay", payreload))
			} else if trxpay == "expire" {
				statustrxcancel := "expired"

				updatetrxcancel := entity.UpdateTransaction(input.Transactions_id, statustrxcancel)
				updatedTrxcancel, err := h.transactionService.UpdateTransactionexp(updatetrxcancel)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}
				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Sorry! payment expired", updatedTrxcancel.Transactions_id))
			} else if trxpay == "cancel" {
				statustrxcancel := "cancel"

				updatetrxcancel := entity.UpdateTransaction(input.Transactions_id, statustrxcancel)
				updatedTrxcancel, err := h.transactionService.UpdateTransactioncancel(updatetrxcancel)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}
				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success cancel", updatedTrxcancel.Transactions_id))
			} else if trxpay == "settlement" {

				statustrx := "settlement"

				updatetrx := entity.UpdateTransaction(input.Transactions_id, statustrx)

				updatedTrx, err := h.transactionService.UpdateTransaction(updatetrx)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}

				transaction_id_detail := uuid.MustParse(input.Transactions_id)
				transactiondetail, err := h.transactionService.FindTrxdetailByID(transaction_id_detail)

				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}

				if transactiondetail.Event_id == "" {
					return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found Event no data"))
				}
				Event_id := uuid.MustParse(transactiondetail.Event_id)
				eventdata, err := h.transactionService.FindEventByID(Event_id)

				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}

				if eventdata.Event_id.String() == "" {
					return c.JSON(http.StatusFound, response.ErrorResponse(http.StatusFound, "Sorry! We found Event no data"))
				}
				trxidfind := uuid.MustParse(transaction.Transactions_id)
				trxdatafind, err := h.transactionService.FindTicketByID(trxidfind)

				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}

				if trxdatafind.Transaction_id == "" {
					ticket_id := uuid.New().String()
					codeqr := uuid.New().String()
					NewTicketdata := entity.NewTicket(ticket_id, transaction.Transactions_id, eventdata.Event_id.String(), codeqr, eventdata.Title_event, transactiondetail.Qty_event)

					ticketdata, err := h.transactionService.CreateTicket(NewTicketdata)
					if err != nil {
						return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
					}

					jsonData, err := json.Marshal(checkpayreq)
					if err != nil {
						return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
					}
					var paymentdata Paymentsdata
					err = json.Unmarshal([]byte(jsonData), &paymentdata)
					if err != nil {
						return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
					}
					paymentid := uuid.MustParse(checkpayreq["transaction_id"].(string))
					paydatafind, err := h.paymentService.FindPayByID(paymentid)

					if err != nil {
						return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
					}

					if paydatafind.Payment_id == "" {
						NewPaymentdata := entity.NewPaymentdata(checkpayreq["transaction_id"].(string), paymentdata.Order_id, paymentdata.Transaction_status, paymentdata.Transaction_time, paymentdata.Settlement_time, paymentdata.Payment_type, paymentdata.Signature_key)

						payment, err := h.paymentService.CreatePayment(NewPaymentdata)
						if err != nil {
							return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
						}
						paydat := make(map[string]interface{})

						paydat["payment_id"] = payment.Payment_id
						paydat["transaksi_id"] = payment.Transaksi_id
						paydat["Status_payment"] = payment.Status_pay
						paydat["Transaksi_time"] = payment.Pay_time
						paydat["Settlement_time"] = payment.Pay_settlement_time
						paydat["Payment_type"] = payment.Pay_type
						paydat["Signature_key"] = payment.Signature_key

						checkpay := make(map[string]interface{})

						checkpay["Ticket_id"] = ticketdata.Tickets_id
						checkpay["Status_Payment"] = checkpayreq["transaction_status"]
						checkpay["Status_Transaksi"] = updatedTrx.Status
						checkpay["Message"] = "Payment Success"
						checkpay["Payment"] = paydat

						return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Check Pay", checkpay))
					} else {
						return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! payment settlement"))
					}

				} else {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! payment settlement"))
				}

			} else {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! System error"))
			}

		} else {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
		}

	} else {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! access denied"))
	}

}

func (h *TransactionHandler) FindAllTransaction(c echo.Context) error {
	var input binder.GetAllRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	if input.Key == "trx" {

		if input.Transactions_id == "" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
		}
		if isValidUUID(input.Transactions_id) {
			trx_id := uuid.MustParse(input.Transactions_id)

			trxdata, err := h.transactionService.FindTrxByID(trx_id)

			if err != nil {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			}
			if trxdata.Transactions_id == "" {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Data Not Found"))
			} else {
				return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success menampilkan data Transaction", trxdata))
			}
		} else {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
		}
	}
	if input.Key == "trx_user" {
		if input.User_id == "" || input.Transactions_id == "" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
		}
		if isValidUUID(input.User_id) {
			if isValidUUID(input.Transactions_id) {
				trx_id := uuid.MustParse(input.Transactions_id)
				User_id := uuid.MustParse(input.User_id)

				trxdatauser, err := h.transactionService.FindTrxrelationByID(trx_id, User_id)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}
				if trxdatauser.Transactions_id == "" && trxdatauser.User_id == "" {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Data Not Found"))
				} else {
					return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success menampilkan data Transaction", trxdatauser))
				}
			} else {
				return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
			}

		} else {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
		}
	}

	if input.Key == "trx_admin_all" {

		if input.User_id == "" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! Colum empty"))
		}
		if isValidUUID(input.User_id) {
			User_id := uuid.MustParse(input.User_id)
			// valuser, err := h.transactionService.FindUserByID(User_id)
			// if err != nil {
			// 	return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
			// }

			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusBadRequest, "you must login first"))
			}
			claimsjwt := user.Claims.(*token.JwtCustomClaims)
			if claimsjwt.Role == "admin" {

				trxdatauser, err := h.transactionService.FindTrxrelationadminByID(User_id)
				if err != nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
				}
				if trxdatauser == nil {
					return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Data Not Found"))
				} else {
					return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success menampilkan data Transaction", trxdatauser))
				}
			}
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! access denied"))
		} else {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Sorry! not UUID"))
		}
	}
	return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Format Key Not Found"))

}

// func (h *TransactionHandler) CancelPayTransaction(c echo.Context) error {

// }
