package handler

import (
	"net/http"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/internal/http/binder"
	"github.com/Kevinmajesta/depublic-backend/internal/service"
	"github.com/Kevinmajesta/depublic-backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return CategoryHandler{categoryService: categoryService}
}

// TODO ADD CATEGORY
// func (h *CategoryHandler) AddCategory(c echo.Context) error {
// 	input := binder.CategoryCreateRequest{}

// 	if err := c.Bind(&input); err != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Wrong input!"))
// 	}

// 	// Cek Existing Category
// 	existingCategory, err := h.categoryService.GetCategoryByName(input.NameCategories)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Error checking existing category"))
// 	}

// 	if existingCategory != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Nama categories sudah ada!"))
// 	}

// 	newCategory := entity.NewCategory(input.NameCategories)

// 	category, err := h.categoryService.AddCategory(newCategory)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
// 	}

// 	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Add to category Success", category))
// }

func (h *CategoryHandler) AddCategory(c echo.Context) error {
	input := binder.CategoryCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Wrong input!"))
	}

	// Check input==nil
	if input.NameCategories == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Name cannot be empty"))
	}

	newCategory := entity.NewCategory(input.NameCategories)

	category, err := h.categoryService.AddCategory(newCategory)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Add to category Success", category))
}

// TODO GET ALL CATEGORY

func (h *CategoryHandler) GetAllCategory(c echo.Context) error {
	category, err := h.categoryService.GetAllCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Get All Category Success!", category))
}

// TODO GET CATEGORY BY ID
func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	category, err := h.categoryService.GetCategoryByID(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found"})
	}

	return c.JSON(http.StatusOK, category)
}

// Category By Param

// TODO GET CATEGORY BY PARAM
// BACKUP
// func (h *CategoryHandler) GetCategoryByParam(c echo.Context) error {
// 	id := c.QueryParam("id")
// 	categoryName := c.QueryParam("name")

// 	if id == "" && categoryName == "" {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID or name parameter is required"})
// 	}

// 	if id != "" {
// 		uuid, err := uuid.Parse(id)
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
// 		}

// 		category, err := h.categoryService.GetCategoryByID(uuid)
// 		if err != nil {
// 			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found!"})
// 		}

// 		return c.JSON(http.StatusOK, category)
// 	}

// 	if categoryName != "" {
// 		category, err := h.categoryService.GetCategoryByName(categoryName)
// 		if err != nil {
// 			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found!"})
// 		}

// 		return c.JSON(http.StatusOK, category)
// 	}

// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error"})
// }

func (h *CategoryHandler) GetCategoryByParam(c echo.Context) error {
	id := c.QueryParam("id")
	categoryName := c.QueryParam("name")

	// Cond if Id and Name Null
	if id == "" && categoryName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID or name is required"})
	}

	// Cond if id not found/wrong Id
	if id != "" {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Id not found"})
		}

		category, err := h.categoryService.GetCategoryByID(uuid)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found!"})
		}
		// Cond if id and name not match
		if categoryName != "" && category.NameCategories != categoryName {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Id and name do not match"})
		}

		return c.JSON(http.StatusOK, category)
	}

	if categoryName != "" {
		category, err := h.categoryService.GetCategoryByName(categoryName)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found!"})
		}
		if category == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found!"})
		}

		return c.JSON(http.StatusOK, category)
	}

	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error"})
}

// TODO DELETE CATEGORY BY ID

func (h *CategoryHandler) DeleteCategoryByID(c echo.Context) error {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	category, err := h.categoryService.DeleteCategoryByID(uuid)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Event category not found"})
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Delete Category Success!", category.DeletedAt))
}

// Update Category

// func (h *CategoryHandler) UpdateCategoryByID(c echo.Context) error {
// 	// categoryID, err := uuid.Parse(c.Param("event_categories_id"))
// 	id := c.Param("id")
// 	uuid, err := uuid.Parse(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
// 	}

// 	input := binder.CategoryUpdateRequest{}
// 	if err := c.Bind(&input); err != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Wrong Input!"))
// 	}

// 	existingCategory, err := h.categoryService.UpdateCategoryByID(uuid, input.NameCategories)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Event Category not found"))
// 	}

// 	existingCategory.NameCategories = input.NameCategories

// 	updatedCategory, err := h.categoryService.UpdateCategoryByID(existingCategory.EventCategoriesID, input.NameCategories)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
// 	}

// 	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Update Category Success!", updatedCategory))
// }

// TRY/ERROR (SUCCESS)
func (h *CategoryHandler) UpdateCategoryByID(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	input := binder.CategoryUpdateRequest{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Wrong Input!"))
	}

	existingCategory, err := h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Category not found"))
	}

	existingCategory.NameCategories = input.NameCategories

	updateCategory, err := h.categoryService.UpdateCategoryByID(existingCategory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Update Category Success!", updateCategory))
}
