package repository

import (
	"encoding/json"
	"time"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistRepository interface {
	GetWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error)
	CheckEventAdd(UserId uuid.UUID) (bool, error)
	AddWishlist(wishlist *entity.Wishlist) (*entity.Wishlist, error)
	GetWishlistByEventAndUser(EventId, UserId uuid.UUID) (*entity.Wishlist, error)
	RemoveWishlist(EventId, UserId uuid.UUID) error
	GetAllWishlist() ([]entity.Wishlist, error)
	FindWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error)
}

type wishlistRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewWishlistRepository(db *gorm.DB, cacheable cache.Cacheable) WishlistRepository {
	return &wishlistRepository{db: db, cacheable: cacheable}
}

func (r *wishlistRepository) GetWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error) {
	wishlist := new(entity.Wishlist)

	if err := r.db.Where("user_id = ?", UserId).Take(wishlist).Error; err != nil {
		return wishlist, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) CheckEventAdd(UserId uuid.UUID) (bool, error) {

	var eventAdd int

	err := r.db.Raw("SELECT COUNT(event_id) FROM wishlists WHERE user_id = ?", UserId).Scan(&eventAdd).Error
	if err != nil {
		return false, err
	}

	return eventAdd == 1, nil
}

func (r *wishlistRepository) AddWishlist(wishlist *entity.Wishlist) (*entity.Wishlist, error) {
	if err := r.db.Create(wishlist).Error; err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) GetWishlistByEventAndUser(EventId, UserId uuid.UUID) (*entity.Wishlist, error) {
	wishlist := new(entity.Wishlist)
	if err := r.db.Where("event_id = ? AND user_id = ?", EventId, UserId).First(wishlist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return wishlist, nil
}

func (r *wishlistRepository) RemoveWishlist(EventId, UserId uuid.UUID) error {

	if err := r.db.Exec("DELETE FROM wishlists WHERE event_id = ? AND user_id = ?", EventId, UserId).Error; err != nil {
		return err
	}

	return nil

}

func (r *wishlistRepository) GetAllWishlist() ([]entity.Wishlist, error) {
	wishlists := make([]entity.Wishlist, 0)

	key := "GetAllWishlists"

	data, _ := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Find(&wishlists).Error; err != nil {
			return wishlists, err
		}
		marshalledWishlists, _ := json.Marshal(wishlists)
		err := r.cacheable.Set(key, marshalledWishlists, 5*time.Minute)
		if err != nil {
			return wishlists, err
		}
	} else {
		// Data ditemukan di Redis, unmarshal data ke users
		err := json.Unmarshal([]byte(data), &wishlists)
		if err != nil {
			return wishlists, err
		}
	}

	return wishlists, nil
}

func (r *wishlistRepository) FindWishlistByUserId(UserId uuid.UUID) (*entity.Wishlist, error) {
	wishlists := new(entity.Wishlist)

	if err := r.db.Raw("SELECT * FROM wishlists WHERE user_id = ?", UserId).Take(wishlists).Error; err != nil {
		return wishlists, err
	}
	return wishlists, nil
}
