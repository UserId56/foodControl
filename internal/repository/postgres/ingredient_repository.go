package postgres

import (
	"context"
	"errors"
	"fmt"
	"foodControl/internal/domain"

	"gorm.io/gorm"
)

type IngredientRepository struct {
	db *gorm.DB
}

type GormIngredient struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"index:name;type:text"`
	Price float32 `gorm:"type:decimal(10,2)"`
	Unit  string  `gorm:"size:15"`
}

func (gi GormIngredient) toDomain() *domain.Ingredient {
	return &domain.Ingredient{
		Id:    gi.ID,
		Name:  gi.Name,
		Price: gi.Price,
		Unit:  gi.Unit,
	}
}

func (gi *GormIngredient) fromDomain(ingredient *domain.Ingredient) {
	gi.ID = ingredient.Id
	gi.Name = ingredient.Name
	gi.Price = ingredient.Price
	gi.Unit = ingredient.Unit
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return IngredientRepository{
		db: db,
	}
}

func (ir IngredientRepository) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	var DBIngredient GormIngredient
	DBIngredient.fromDomain(ingredient)
	if err := ir.db.WithContext(ctx).Create(&DBIngredient).Error; err != nil {
		fmt.Println("Ошибка создания ингридиента: ", err)
		return domain.ErrCreate
	}
	*ingredient = *DBIngredient.toDomain()
	return nil
}

func (ir IngredientRepository) CheckName(ctx context.Context, ingredient *domain.Ingredient) error {
	var DBIngredient GormIngredient
	if err := ir.db.WithContext(ctx).Select("ID").Where("name = ?", ingredient.Name).First(&DBIngredient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		fmt.Println("Ошибка проверки ингридиента на дубль: ", err)
		return domain.ErrInternal
	}
	if DBIngredient.ID != 0 {
		return domain.ErrIngredientAlreadyExists
	}
	return nil
}
