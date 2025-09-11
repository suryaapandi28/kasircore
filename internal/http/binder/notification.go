package binder

type CreateNotification struct {
	UserId  string `json:"user_id" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Message string `json:"message" validate:"required"`
	Is_Read bool   `json:"is_read"`
}

type MarkNotificationAsRead struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetAllRequestNotif struct {
	Key    string `json:"key" validate:"required,key"`
	UserId string `json:"user_id" validate:"required"`
}
