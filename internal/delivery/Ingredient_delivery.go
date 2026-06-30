package delivery

import (
	"context"
	"errors"
	"fmt"
	"foodControl/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IngredientUseCase interface {
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	// Update(ctx context.Context, ingredient *domain.Ingredient) error
	// Delete(ctx context.Context, ingredient *domain.Ingredient) error
	// GetByIds(ctx context.Context, ids []uint) ([]*domain.Ingredient, error)
}

type IngredientHandler struct {
	IngredientUseCase IngredientUseCase
}

func RegisterEndpoints(router *gin.Engine, uc IngredientUseCase) {
	handler := &IngredientHandler{
		IngredientUseCase: uc,
	}
	ingredient := router.Group("/ingredient")
	ingredient.POST("/create", handler.Create)
}

type IngredientCreate struct {
	Name  string  `json:"name" binding:"required"`
	Price float32 `json:"price" binding:"required"`
	Unit  string  `json:"unit" binding:"required"`
}

func (ic IngredientCreate) toDomain() *domain.Ingredient {
	return &domain.Ingredient{
		Name:  ic.Name,
		Price: ic.Price,
		Unit:  ic.Unit,
	}
}

func (ih IngredientHandler) Create(c *gin.Context) {
	var req IngredientCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Ошибка парсинга тела: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    fmt.Sprintf("Ошибка парсинга тела: %+v", err),
		})
		return
	}
	ingredient := req.toDomain()
	if err := ih.IngredientUseCase.Create(c.Request.Context(), ingredient); err != nil {
		if errors.Is(err, domain.ErrIngredientAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    err.Error(),
			})
			return
		}
		fmt.Println("Ошибка при создании ингридиента: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    domain.ErrInternal,
		})
		return
	}
	c.JSON(http.StatusCreated, ingredient)
}
