package handler

import (
	"log"
	"net/http"

	"github.com/Kevinmajesta/depublic-backend/internal/service"
	"github.com/Kevinmajesta/depublic-backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) TicketHandler {
	if ticketService == nil {
		log.Fatal("ticketService must not be nil")
	}
	return TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) FindAllTicket(c echo.Context) error {
	tickets, err := h.ticketService.FindAllTicket()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays ticket data", tickets))
}

func (h *TicketHandler) FindTicketsByEventID(c echo.Context) error {
	eventIDParam := c.Param("eventID")
	eventUUID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid event ID"))
	}
	exists, err := h.ticketService.CheckTicketExists(eventUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "event ID does not exist"))
	}
	tickets, err := h.ticketService.FindTicketsByEventID(eventUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	if tickets == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "event ID not found"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays ticket data for the event", tickets))
}

func (h *TicketHandler) FindTicketsByQRCode(c echo.Context) error {
	QRCodeParam := c.Param("QRCode")
	QRCodeUUID, err := uuid.Parse(QRCodeParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid QRCode"))
	}
	exists, err := h.ticketService.CheckTicketExists(QRCodeUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "event ID does not exist"))
	}
	qrcodes, err := h.ticketService.FindTicketsByQRCode(QRCodeUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays ticket data for the QRCode", qrcodes))
}
