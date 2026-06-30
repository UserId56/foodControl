package main

import (
	"fmt"
	"foodControl/internal/delivery"
	"foodControl/internal/usecase"
	"log"
	"os"

	repoPostgres "foodControl/internal/repository/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Start")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=foodcontrol port=5432 sslmode=disable", DBUser, DBPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	db.AutoMigrate(&repoPostgres.GormIngredient{})
	ingredientRepo := repoPostgres.NewIngredientRepository(db)
	ingredientUC := usecase.NewIngredientUseCase(ingredientRepo)
	router := gin.Default()
	delivery.RegisterEndpoints(router, ingredientUC)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
