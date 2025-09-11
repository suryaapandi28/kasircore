package binder

type FindWishlistByUserIdRequest struct {
	UserId string `param:"id" validate:"required"`
}

type WishlistRequest struct {
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`
}

type RemoveWishlistRequest struct {
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`
}
