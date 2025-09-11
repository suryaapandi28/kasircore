package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type NotificationHandler struct {
	notificationService service.NotificationService
	userService         service.UserService
}

func NewNotificationHandler(notificationService service.NotificationService, userService service.UserService) NotificationHandler {
	return NotificationHandler{notificationService: notificationService, userService: userService}
}

func (h *NotificationHandler) GetUserNotifications(c echo.Context) error {
	input := binder.MarkNotificationAsRead{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	userID, err := uuid.Parse(input.UserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid user ID"))
	}

	exists, err := h.notificationService.CheckUserExists(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "user ID does not exist"))
	}

	notifications, err := h.notificationService.GetUserNotifications(userID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	if len(notifications) == 0 {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "No notifications found", notifications))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Get All Notifications", notifications))

}

func (h *NotificationHandler) GetUserNotificationNoRead(c echo.Context) error {
	input := binder.MarkNotificationAsRead{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	userID, err := uuid.Parse(input.UserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid user ID"))
	}

	exists, err := h.notificationService.CheckUserExists(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "user ID does not exist"))
	}

	notifications, err := h.notificationService.GetUserNotificationsNoRead(userID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse(http.StatusUnprocessableEntity, err.Error()))
	}

	if len(notifications) == 0 {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "No notifications found", notifications))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Success Get All Notifications", notifications))
}

// CreateNotification
func (h *NotificationHandler) CreateNotification(c echo.Context) error {
	var input entity.Notification

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	// Set the type and message for the notification, assuming they come from the input
	notification := &entity.Notification{
		Type:    input.Type,
		Message: input.Message,
		IsRead:  false,
	}

	if input.Type == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Notification type cannot be empty"))
	}
	if input.Message == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Notification message cannot be empty"))
	}

	// Call the service to create the notification for all users
	if err := h.notificationService.CreateNotification(notification); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusCreated, "Notification created successfully", nil))
}

// MarkNotificationAsRead
func (h *NotificationHandler) MarkNotificationAsRead(c echo.Context) error {
	notificationID, err := uuid.Parse(c.Param("notification_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid notification ID"))
	}

	if err := h.notificationService.MarkNotificationAsRead(notificationID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}
