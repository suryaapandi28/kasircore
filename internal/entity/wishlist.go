package entity

import "github.com/google/uuid"

type Wishlist struct {
	WishlistId uuid.UUID `json:"wishlist_id" gorm:"primarykey"`
	UserId     uuid.UUID `json:"user_id" gorm:"not null"`
	User       User      `json:"-" gorm:"foreignkey:UserId"`
	EventId    uuid.UUID `json:"event_id" gorm:"not null"`
	Event      Events    `json:"-" gorm:"foreignkey:EventId;constraint:OnDelete:CASCADE"`
	Auditable
}

func NewWishlist(UserId, EventId uuid.UUID) *Wishlist {
	return &Wishlist{
		WishlistId: uuid.New(),
		UserId:     UserId,
		EventId:    EventId,
		Auditable:  NewAuditable(),
	}
}
