package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"kasir-api/internal/config"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/service"
	"kasir-api/pkg/database"

	_ "kasir-api/docs" // Swagger docs

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Kasir API
// @version 1.0
// @description API for Point of Sale (Kasir) system with product and category management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@kasir-api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database
	dbCfg := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	}
	db, err := database.NewPostgres(dbCfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Run Migrations (for simplicity, running it here)
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Dependency Injection Wiring
	// Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)

	// Route Registration
	mux := http.NewServeMux()

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Categories
	mux.HandleFunc("/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)

	// Products
	mux.HandleFunc("/products", productHandler.HandleProducts)
	mux.HandleFunc("/products/", productHandler.HandleProductByID)

	fmt.Println("Server starting on port 8080...")
	fmt.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	);`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price INT NOT NULL,
		stock INT NOT NULL,
		category_id INT REFERENCES categories(id)
	);`

	if _, err := db.Exec(createCategoriesTable); err != nil {
		return fmt.Errorf("error creating categories table: %w", err)
	}

	if _, err := db.Exec(createProductsTable); err != nil {
		return fmt.Errorf("error creating products table: %w", err)
	}

	return nil
}
