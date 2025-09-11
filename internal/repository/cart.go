package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetAllCart() ([]entity.Carts, error)
	FindCartById(CartId uuid.UUID) (*entity.Carts, error)
	GetCartByUserId(UserId uuid.UUID) (*entity.Carts, error)
	GetCartByUserAndEvent(UserId, EventId uuid.UUID) (*entity.Carts, error)
	CheckIfEventAlreadyAdded(UserId, EventId uuid.UUID) (bool, error)
	CheckEventAdd(UserId uuid.UUID) (bool, error)
	CreateCart(cart *entity.Carts) error
	UpdateQuantityAdd(UserId, EventId uuid.UUID) error
	UpdateQuantityLess(UserId, EventId uuid.UUID) error
	UpdateTotalPrice(UserId, EventId uuid.UUID, price int) error
	GetUserTotalQtyInCart(UserId, EventId uuid.UUID) (int, error)
	RemoveCart(cart *entity.Carts) (bool, error)
}
type cartRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewCartRepository(db *gorm.DB, cacheable cache.Cacheable) CartRepository {
	return &cartRepository{db: db, cacheable: cacheable}
}

func (r *cartRepository) GetAllCart() ([]entity.Carts, error) {
	carts := make([]entity.Carts, 0)

	key := "GetAllCarts"

	data, _ := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Find(&carts).Error; err != nil {
			return carts, err
		}
		marshalledCarts, _ := json.Marshal(carts)
		err := r.cacheable.Set(key, marshalledCarts, 1*time.Minute)
		if err != nil {
			return carts, err
		}
	} else {
		// Data ditemukan di Redis, unmarshal data ke users
		err := json.Unmarshal([]byte(data), &carts)
		if err != nil {
			return carts, err
		}
	}

	return carts, nil
}

func (r *cartRepository) FindCartById(CartId uuid.UUID) (*entity.Carts, error) {
	cart := &entity.Carts{}

	if err := r.db.Where("cart_id = ?", CartId).Take(cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (r *cartRepository) GetCartByUserId(UserId uuid.UUID) (*entity.Carts, error) {
	cart := new(entity.Carts)

	if err := r.db.Where("user_id = ?", UserId).Take(cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

func (r *cartRepository) GetCartByUserAndEvent(UserId, EventId uuid.UUID) (*entity.Carts, error) {
	var cart entity.Carts

	err := r.db.Raw("SELECT * FROM carts WHERE user_id = ? AND event_id = ?", UserId, EventId).Scan(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Tidak ada record ditemukan
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) CheckIfEventAlreadyAdded(UserId, EventId uuid.UUID) (bool, error) {

	var qty int

	if err := r.db.Raw("SELECT qty FROM carts WHERE user_id = $1 AND event_id = $2", UserId, EventId).Scan(&qty).Error; err != nil {
		return false, err
	}

	return qty >= 1, nil
}

func (r *cartRepository) CheckEventAdd(UserId uuid.UUID) (bool, error) {

	var eventAdd int

	err := r.db.Raw("SELECT COUNT(event_id) FROM carts WHERE user_id = ?", UserId).Scan(&eventAdd).Error
	if err != nil {
		return false, err
	}

	return eventAdd >= 1, nil
}

func (r *cartRepository) CreateCart(cart *entity.Carts) error {
	err := r.db.Exec(
		"INSERT INTO carts (cart_id, user_id, event_id, qty, ticket_date, price) VALUES (?, ?, ?, ?, ?, ?)",
		cart.Cart_id, cart.User_id, cart.Event_id, cart.Qty, cart.Ticket_date, cart.Price,
	).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) UpdateQuantityAdd(UserId, EventId uuid.UUID) error {
	if err := r.db.Exec("UPDATE carts SET qty = qty + 1 WHERE user_id = ? AND event_id = ? AND qty > 0", UserId, EventId).Error; err != nil {
		return err
	}

	return nil
}
func (r *cartRepository) UpdateQuantityLess(UserId, EventId uuid.UUID) error {
	if err := r.db.Exec("UPDATE carts SET qty = qty - 1 WHERE user_id = ? AND event_id = ? AND qty > 0", UserId, EventId).Error; err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) UpdateTotalPrice(UserId, EventId uuid.UUID, price int) error {
	query := `
        UPDATE carts
        SET price = $1
        WHERE user_id = $2 AND event_id = $3
    `
	err := r.db.Exec(query, price, UserId, EventId).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *cartRepository) GetUserTotalQtyInCart(UserId, EventId uuid.UUID) (int, error) {
	var totalQty int
	err := r.db.Raw(
		"SELECT COALESCE(SUM(qty), 0) AS total_qty FROM carts WHERE user_id = ? AND event_id = ?",
		UserId, EventId,
	).Scan(&totalQty).Error
	if err != nil {
		return 0, err
	}
	return totalQty, nil
}

func (r *cartRepository) RemoveCart(cart *entity.Carts) (bool, error) {
	// Hapus entri dari keranjang (hard delete)
	if err := r.db.Unscoped().Delete(cart).Error; err != nil {
		return false, err
	}
	return true, nil
}
