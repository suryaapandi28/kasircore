package repository

import (
	"encoding/json"
	"time"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindEventByID(Event_id uuid.UUID) (*entity.Events, error)
	FindTrxByID(Transaction_id uuid.UUID) (*entity.Transactions, error)
	FindTrxrelationByID(Transaction_id uuid.UUID, User_id uuid.UUID) (*entity.Transactions, error)
	FindTrxrelationadminByID(User_id uuid.UUID) ([]entity.Transactions, error)
	FindTrxdetailByID(Transaction_id uuid.UUID) (*entity.Transaction_details, error)
	FindCartByID(Cart_id uuid.UUID) (*entity.Carts, error)
	CreateTransaction(transaction *entity.Transactions) (*entity.Transactions, error)
	CreateTransactiondetail(transaction *entity.Transaction_details) (*entity.Transaction_details, error)
	CreateTicket(transaction *entity.Tickets) (*entity.Tickets, error)
	FindAllTransaction() ([]entity.Transactions, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindUserByID(cart_id uuid.UUID) (*entity.Useraccount, error)
	UpdateTransaction(transactionupdate *entity.Transactions) (*entity.Transactions, error)
	UpdateTransactioncancel(transactionupdate *entity.Transactions) (*entity.Transactions, error)
	UpdateTransactionexp(transactionupdate *entity.Transactions) (*entity.Transactions, error)
	FindTicketByID(Transaction_id uuid.UUID) (*entity.Tickets, error)
}

type transactionRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

type Result struct {
	// Define fields according to your SQL query result columns
}

func NewTransactionRepository(db *gorm.DB, cacheable cache.Cacheable) TransactionRepository {
	return &transactionRepository{db: db, cacheable: cacheable}
}

func (r *transactionRepository) CreateTransaction(transaction *entity.Transactions) (*entity.Transactions, error) {

	if err := r.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *transactionRepository) CreateTransactiondetail(transactiondetail *entity.Transaction_details) (*entity.Transaction_details, error) {

	if err := r.db.Create(&transactiondetail).Error; err != nil {
		return transactiondetail, err
	}
	return transactiondetail, nil

}

func (r *transactionRepository) CreateTicket(ticket *entity.Tickets) (*entity.Tickets, error) {

	if err := r.db.Create(&ticket).Error; err != nil {
		return ticket, err
	}
	return ticket, nil

}

func (r *transactionRepository) FindAllTransaction() ([]entity.Transactions, error) {
	transaction := make([]entity.Transactions, 0)

	key := "FindAllTransactions"

	data, _ := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Find(&transaction).Error; err != nil {
			return transaction, err
		}
		marshalledtransaction, _ := json.Marshal(transaction)
		err := r.cacheable.Set(key, marshalledtransaction, 5*time.Minute)
		if err != nil {
			return transaction, err
		}
	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &transaction)
		if err != nil {
			return transaction, err
		}
	}

	return transaction, nil
}
func (r *transactionRepository) FindEventByID(Event_id uuid.UUID) (*entity.Events, error) {
	var events entity.Events

	eventsdata := &events
	key := "FindEventByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw("SELECT * FROM events WHERE event_id = ?", Event_id).Scan(&events).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &events, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &eventsdata)
		if err != nil {
			return eventsdata, err
		}
	}
	return &events, nil

}

func (r *transactionRepository) FindTrxByID(Transaction_id uuid.UUID) (*entity.Transactions, error) {
	var trx entity.Transactions

	trxsdata := &trx
	key := "FindtrxByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw("SELECT * FROM transactions WHERE transactions_id = ?", Transaction_id).Scan(&trx).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &trx, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &trxsdata)
		if err != nil {
			return trxsdata, err
		}
	}
	return &trx, nil

}

func (r *transactionRepository) FindTicketByID(Transaction_id uuid.UUID) (*entity.Tickets, error) {
	var ticket entity.Tickets

	ticketsdata := &ticket
	key := "FindTicketByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw("SELECT * FROM tickets WHERE transaction_id = ?", Transaction_id).Scan(&ticket).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &ticket, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &ticketsdata)
		if err != nil {
			return ticketsdata, err
		}
	}
	return &ticket, nil

}

func (r *transactionRepository) FindTrxrelationByID(Transaction_id uuid.UUID, User_id uuid.UUID) (*entity.Transactions, error) {
	var trx entity.Transactions

	trxsdata := &trx
	key := "FindTrxrelationByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw(`SELECT t.* FROM transactions t JOIN users u ON t.user_id = ? AND t.transactions_id = ?`, User_id, Transaction_id).Scan(&trx).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &trx, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &trxsdata)
		if err != nil {
			return trxsdata, err
		}
	}
	return &trx, nil

}

func (r *transactionRepository) FindTrxrelationadminByID(User_id uuid.UUID) ([]entity.Transactions, error) {

	var trx []entity.Transactions

	// trxsdata := &trx
	key := "FindTrxrelationByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw(`SELECT t.* FROM transactions t JOIN users u ON t.status = true`).Scan(&trx).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return trx, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &trx)
		if err != nil {
			return trx, err
		}
	}
	return trx, nil

	// var transactions []entity.Transactions

	// key := "FindAllTransactions"

	// data, _ := r.cacheable.Get(key)
	// // roleadmin := "admin"
	// if data == "" {
	// 	if err := r.db.Table("transactions").
	// 		Select("transactions.*").
	// 		Joins("INNER JOIN users ON transactions.user_id = ?", User_id).
	// 		Where("users.role = ?", "admin").
	// 		Find(&transactions).Error; err != nil {
	// 		return transactions, err
	// 	}
	// 	marshalledtransaction, _ := json.Marshal(transactions)
	// 	err := r.cacheable.Set(key, marshalledtransaction, 5*time.Minute)
	// 	if err != nil {
	// 		return transactions, err
	// 	}
	// } else {
	// 	// Data ditemukan di Redis, unmarshal data ke transaction
	// 	err := json.Unmarshal([]byte(data), &transactions)
	// 	if err != nil {
	// 		return transactions, err
	// 	}
	// }

	// return transactions, nil

	// transaction := make([]entity.Transactions, 0)

	// trxsdata := &transaction
	// // roleadm := "admin"
	// err := r.db.Table("transactions").
	// 	Select("transactions.*").
	// 	Joins("INNER JOIN users ON transactions.user_id = ?", User_id).
	// 	// Where("users.role = ?", "admin").
	// 	Find(&transaction).Error

	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }
	// return trxsdata, nil
}

func (r *transactionRepository) FindTrxdetailByID(Transaction_id uuid.UUID) (*entity.Transaction_details, error) {
	var trx entity.Transaction_details

	trxsdata := &trx
	key := "FindtrxdetailByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw("SELECT * FROM transaction_details WHERE transaction_id = ?", Transaction_id).Scan(&trx).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &trx, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &trxsdata)
		if err != nil {
			return trxsdata, err
		}
	}
	return &trx, nil

}

func (r *transactionRepository) FindCartByID(Cart_id uuid.UUID) (*entity.Carts, error) {

	var cart entity.Carts

	cartdata := &cart
	key := "FindEventByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		var cartdata entity.Carts
		err := r.db.Raw("SELECT * FROM carts WHERE cart_id = ?", Cart_id).Scan(&cartdata).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &cartdata, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &cartdata)
		if err != nil {
			return cartdata, err
		}
	}
	return &cart, nil
}

func (r *transactionRepository) FindUserByID(User_ID uuid.UUID) (*entity.Useraccount, error) {

	var user entity.Useraccount

	userdata := &user
	key := "FindEventByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		var userdata entity.Useraccount
		err := r.db.Raw("SELECT * FROM users WHERE user_id = ?", User_ID).Scan(&userdata).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &userdata, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &userdata)
		if err != nil {
			return userdata, err
		}
	}
	return &user, nil
}

func (r *transactionRepository) FindUserByEmail(email string) (*entity.User, error) {
	user := new(entity.User)
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *transactionRepository) UpdateTransaction(transactionupdate *entity.Transactions) (*entity.Transactions, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if transactionupdate.Status != "" {
		fields["status"] = transactionupdate.Status
	}

	// Update the database in one query.
	if err := r.db.Model(transactionupdate).Where("transactions_id = ?", transactionupdate.Transactions_id).Updates(fields).Error; err != nil {
		return transactionupdate, err
	}

	return transactionupdate, nil
}

func (r *transactionRepository) UpdateTransactioncancel(transactionupdate *entity.Transactions) (*entity.Transactions, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if transactionupdate.Status != "" {
		fields["status"] = transactionupdate.Status
	}

	// Update the database in one query.
	if err := r.db.Model(transactionupdate).Where("transactions_id = ?", transactionupdate.Transactions_id).Updates(fields).Error; err != nil {
		return transactionupdate, err
	}

	return transactionupdate, nil
}
func (r *transactionRepository) UpdateTransactionexp(transactionupdate *entity.Transactions) (*entity.Transactions, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if transactionupdate.Status != "" {
		fields["status"] = transactionupdate.Status
	}

	// Update the database in one query.
	if err := r.db.Model(transactionupdate).Where("transactions_id = ?", transactionupdate.Transactions_id).Updates(fields).Error; err != nil {
		return transactionupdate, err
	}

	return transactionupdate, nil
}

// func (r *transactionRepository) UpdateTransaction(data *entity.Transactions) (*entity.Transactions, error) {
// 	err := r.db.Model(&entity.Transactions{}).Where("id = ?", data.Transactions_id).Updates(data).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
