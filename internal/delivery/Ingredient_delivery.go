package delivery

import (
	"context"
	"errors"
	"fmt"
	"foodControl/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IngredientUseCase interface {
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	GetById(ctx context.Context, ids uint) (*domain.Ingredient, error)
	Update(ctx context.Context, ingredient *domain.Ingredient) error
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
	ingredient.GET("/:id", handler.GetById)
	ingredient.POST("/edit", handler.Update)
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
		if errors.Is(err, domain.ErrNameDublicate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    domain.ErrInternal.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, ingredient)
}

func (ih IngredientHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	parUintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println("Ошибка парсинга id ингридиента: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "Не корректный id",
		})
		return
	}
	uintId := uint(parUintId)
	ingridint, err := ih.IngredientUseCase.GetById(c.Request.Context(), uintId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"msg":    domain.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ingridint)
}

func (ih IngredientHandler) Update(c *gin.Context) {
	var req domain.Ingredient
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Ошибка парсинга тела: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    fmt.Sprintf("Ошибка парсинга тела: %+v", err),
		})
		return
	}
	if err := ih.IngredientUseCase.Update(c.Request.Context(), &req); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"msg":    domain.ErrNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, req)
}
