package postgres

import (
	"context"
	"errors"
	"fmt"
	"foodControl/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type IngredientRepository struct {
	db *gorm.DB
}

type GormIngredient struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"index:name;type:text;unique"`
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// SQLSTATE 23505 = unique_violation
			if pgErr.Code == "23505" {
				return domain.ErrNameDublicate
			}
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrNameDublicate
		}
		fmt.Println("Ошибка создания ингридиента: ", err)
		return domain.ErrCreate
	}
	*ingredient = *DBIngredient.toDomain()
	return nil
}

func (ir IngredientRepository) GetById(ctx context.Context, id uint) (*domain.Ingredient, error) {
	var gIngredient GormIngredient
	if err := ir.db.WithContext(ctx).Where("ID = ?", id).First(&gIngredient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		fmt.Println("Ошибка получения из БД: ", err)
		return nil, domain.ErrInternal
	}
	return gIngredient.toDomain(), nil
}

func (ir IngredientRepository) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	var DBIngredient GormIngredient
	DBIngredient.fromDomain(ingredient)
	result := ir.db.WithContext(ctx).Model(&GormIngredient{}).Where("ID = ?", ingredient.Id).Updates(&DBIngredient)
	if result.Error != nil {
		fmt.Println("Ошибка создания ингридиента: ", result.Error)
		return domain.ErrInternal
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	*ingredient = *DBIngredient.toDomain()
	return nil
}
