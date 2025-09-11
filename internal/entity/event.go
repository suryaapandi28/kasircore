package entity

import (
	"time"

	"github.com/google/uuid"
)

//	type Events struct {
//		Event_id          string `json:"event_id"  gorm:"primarykey"`
//		Category_id       string `json:"category_id"`
//		Title_event       string `json:"title_event"`
//		Date_event        string `json:"date_event"`
//		Price_event       string `json:"price_event"`
//		City_event        string `json:"city_event"`
//		Address_event     string `json:"address_event"`
//		Qty_event         string `json:"qty_event"`
//		Description_event string `json:"description_event"`
//		Image_url         string `json:"image_url"`
//		Auditable
//	}
type Events struct {
	Event_id          uuid.UUID `json:"event_id" gorm:"type:uuid;primary_key"`
	Category_id       uuid.UUID `json:"category_id"`
	Title_event       string    `json:"title_event"`
	Date_event        string    `json:"date_event"`
	Price_event       int       `json:"price_event"`
	City_event        string    `json:"city_event"`
	Address_event     string    `json:"address_event"`
	Qty_event         int       `json:"qty_event"`
	Description_event string    `json:"description_event"`
	Image_url         string    `json:"image_url"`
	Auditable
}

func NewEvent(
	categoryID uuid.UUID,
	titleEvent string,
	dateEvent string,
	priceEvent int,
	cityEvent string,
	addressEvent string,
	qtyEvent int,
	descriptionEvent string,
	imageURL string,
) *Events {
	return &Events{
		Event_id:          uuid.New(),
		Category_id:       categoryID,
		Title_event:       titleEvent,
		Date_event:        dateEvent,
		Price_event:       priceEvent,
		City_event:        cityEvent,
		Address_event:     addressEvent,
		Qty_event:         qtyEvent,
		Description_event: descriptionEvent,
		Image_url:         imageURL,
		Auditable:         NewAuditable(),
	}
}

func UpdateEvent(
	event *Events,
	categoryID uuid.UUID,
	titleEvent string,
	dateEvent time.Time,
	priceEvent int,
	cityEvent string,
	addressEvent string,
	qtyEvent int,
	descriptionEvent string,
	imageURL string,
) *Events {
	event.Category_id = categoryID
	event.Title_event = titleEvent
	event.Date_event = dateEvent.Format("2000-01-01")
	event.Price_event = priceEvent
	event.City_event = cityEvent
	event.Address_event = addressEvent
	event.Qty_event = qtyEvent
	event.Description_event = descriptionEvent
	if imageURL != "" {
		event.Image_url = imageURL
	}
	event.Auditable = UpdateAuditable()
	return event
}
