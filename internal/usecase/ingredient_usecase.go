package usecase

import (
	"context"
	"foodControl/internal/domain"
)

type IngredientRepository interface {
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	CheckName(ctx context.Context, ingredient *domain.Ingredient) error
	// Update(ctx context.Context, ingredient *domain.Ingredient) error
	// Delete(ctx context.Context, ingredient *domain.Ingredient) error
	// GetByIds(ctx context.Context, ids []uint) ([]*domain.Ingredient, error)
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
	if err := iuc.IngredientRepository.CheckName(ctx, ingredient); err != nil {
		return err
	}
	if err := iuc.IngredientRepository.Create(ctx, ingredient); err != nil {
		return err
	}
	return nil
}
