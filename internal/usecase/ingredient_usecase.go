package usecase

import (
	"context"
	"foodControl/internal/domain"
)

type IngredientRepository interface {
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	GetById(ctx context.Context, id uint) (*domain.Ingredient, error)
	Update(ctx context.Context, ingredient *domain.Ingredient) error
	// Delete(ctx context.Context, ingredient *domain.Ingredient) error
	GetByIds(ctx context.Context, ids []uint) ([]*domain.Ingredient, error)
}

type IngredientUseCase struct {
	IngredientRepository IngredientRepository
}

func NewIngredientUseCase(repo IngredientRepository) IngredientUseCase {
	return IngredientUseCase{
		IngredientRepository: repo,
	}
}

func (iuc IngredientUseCase) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	if err := iuc.IngredientRepository.Create(ctx, ingredient); err != nil {
		return err
	}
	return nil
}

func (iuc IngredientUseCase) GetById(ctx context.Context, id uint) (*domain.Ingredient, error) {
	ingridient, err := iuc.IngredientRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return ingridient, nil
}

func (iuc IngredientUseCase) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	if err := iuc.IngredientRepository.Update(ctx, ingredient); err != nil {
		return err
	}
	return nil
}

func (iuc IngredientUseCase) GetByIds(ctx context.Context, ids []uint) ([]*domain.Ingredient, error) {
	ingridients, err := iuc.IngredientRepository.GetByIds(ctx, ids)
	if err != nil {
		return ingridients, err
	}
	return ingridients, nil
}
