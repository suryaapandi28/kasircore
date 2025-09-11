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

type WishlistHandler struct {
	wishlistService service.WishlistService
}

func NewWishlistHandler(wishlistService service.WishlistService) WishlistHandler {
	return WishlistHandler{wishlistService: wishlistService}
}

func (h *WishlistHandler) GetWishlistByUserId(c echo.Context) error {
	input := binder.FindWishlistByUserIdRequest{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	userId := uuid.MustParse(input.UserId)

	wishlist, err := h.wishlistService.GetWishlistByUserId(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully showed wishlist data by users", wishlist))
}

func (h *WishlistHandler) GetAllWishlist(c echo.Context) error {
	wishlists, err := h.wishlistService.GetAllWishlist()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully retrieves wishlist data", wishlists))
}

func (h *WishlistHandler) AddWishlist(c echo.Context) error {
	input := binder.WishlistRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	userId := uuid.MustParse(input.UserId)
	eventId := uuid.MustParse(input.EventId)

	newWishlist := entity.NewWishlist(userId, eventId)

	wishlist, err := h.wishlistService.AddWishlist(newWishlist)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses menambahkan wishlist", wishlist))
}

func (h *WishlistHandler) RemoveWishlist(c echo.Context) error {
	var input binder.RemoveWishlistRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	eventId := uuid.MustParse(input.EventId)
	userId := uuid.MustParse(input.UserId)

	wishlistID, err := h.wishlistService.RemoveWishlist(eventId, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success remove wishlist", wishlistID))
}
