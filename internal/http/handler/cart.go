package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return CartHandler{cartService: cartService}
}

func (h *CartHandler) GetAllCarts(c echo.Context) error {
	carts, err := h.cartService.GetAllCart()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully retrieves cart data", carts))
}

func (h *CartHandler) AddToCart(c echo.Context) error {
	input := binder.AddCartRequest{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	eventId := uuid.MustParse(input.EventId)
	userId := uuid.MustParse(input.UserId)

	// Memanggil service untuk menambahkan ke keranjang
	cart, err := h.cartService.AddToCart(userId, eventId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Jika tidak ada kesalahan, kembalikan data AddCartResponse dalam respon JSON
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully added to car", cart))
}

func (h *CartHandler) UpdateQuantityAdd(c echo.Context) error {
	input := binder.UpdateQuantityLessRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "there is an input error",
			},
		})
	}

	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid user_id format",
			},
		})
	}

	eventId, err := uuid.Parse(input.EventId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid event_id format",
			},
		})
	}

	// Panggil service untuk mengurangi quantity dalam keranjang
	err = h.cartService.UpdateQuantityAdd(userId, eventId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"meta": map[string]interface{}{
			"code":    http.StatusOK,
			"message": "successfully add quantity in cart",
		},
	})
}

func (h *CartHandler) UpdateQuantityLess(c echo.Context) error {
	input := binder.UpdateQuantityLessRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "there is an input error",
			},
		})
	}

	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid user_id format",
			},
		})
	}

	eventId, err := uuid.Parse(input.EventId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid event_id format",
			},
		})
	}

	// Panggil service untuk mengurangi quantity dalam keranjang
	err = h.cartService.UpdateQuantityLess(userId, eventId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"meta": map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"meta": map[string]interface{}{
			"code":    http.StatusOK,
			"message": "successfully less quantity in cart",
		},
	})
}

func (h *CartHandler) GetCartByUserId(c echo.Context) error {
	input := binder.FindCartByUserIdRequest{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	userId := uuid.MustParse(input.UserID)

	// Memanggil service untuk menambahkan ke keranjang
	cart, err := h.cartService.GetCartByUserId(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Jika tidak ada kesalahan, kembalikan data AddCartResponse dalam respon JSON
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully showed cart data by users", cart))
}

func (h *CartHandler) RemoveCart(c echo.Context) error {
	var input binder.RemoveCartRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	cartId := uuid.MustParse(input.CartID)

	CartID, err := h.cartService.RemoveCart(cartId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully removed from cart", CartID))
}
