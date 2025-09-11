package binder

import (
	"mime/multipart"

	"github.com/google/uuid"
)

// TODO Add event with image

type EventAddRequest struct {
	CategoryID       uuid.UUID             `form:"category_id" validate:"required"`
	TitleEvent       string                `form:"title_event" validate:"required"`
	DateEvent        string                `form:"date_event" validate:"required"`
	PriceEvent       int                   `form:"price_event" validate:"required"`
	CityEvent        string                `form:"city_event" validate:"required"`
	AddressEvent     string                `form:"address_event" validate:"required"`
	QtyEvent         int                   `form:"qty_event" validate:"required"`
	DescriptionEvent string                `form:"description_event" validate:"required"`
	Image            *multipart.FileHeader `form:"image" validate:"required"`
}
type EventUpdateRequest struct {
	CategoryID       uuid.UUID             `form:"category_id" validate:"required"`
	TitleEvent       string                `form:"title_event" validate:"required"`
	DateEvent        string                `form:"date_event" validate:"required"`
	PriceEvent       int                   `form:"price_event" validate:"required"`
	CityEvent        string                `form:"city_event" validate:"required"`
	AddressEvent     string                `form:"address_event" validate:"required"`
	QtyEvent         int                   `form:"qty_event" validate:"required"`
	DescriptionEvent string                `form:"description_event" validate:"required"`
	Image            *multipart.FileHeader `form:"image"  validate:"required"`
}
