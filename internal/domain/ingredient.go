package domain

import "errors"

type Ingredient struct {
	Id    uint
	Name  string
	Price float32
	Unit  string
}

var (
	ErrInternal                = errors.New("Ошибка на сервере")
	ErrIngredientAlreadyExists = errors.New("ингредиент с таким именем существует.")
	ErrCreate                  = errors.New("ошибка создания элемента")
	ErrNotFound                = errors.New("Элемент не найден")
)
