package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"

	"kasir-api-irfan/database"
	"kasir-api-irfan/docs"
	_ "kasir-api-irfan/docs"
	"kasir-api-irfan/handlers"
	"kasir-api-irfan/repositories"
	"kasir-api-irfan/services"
	"log"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Health check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "API is healthy",
	})
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// @title Kasir API
// @version 1.0
// @BasePath /
func main() {
	// Load configuration
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Create a mux with CORS middleware
	mux := http.NewServeMux()

	// check api health
	mux.HandleFunc("/health", HealthCheck)

	// swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// scalar documentation
	mux.HandleFunc("/scalar", func(w http.ResponseWriter, r *http.Request) {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: docs.SwaggerInfo.ReadDoc(),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Kasir API Reference",
			},
		})
		if err != nil {
			log.Printf("Error generating Scalar HTML: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	// Apply CORS middleware to all routes
	handler := corsMiddleware(mux)

	// Initialize Database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("CRITICAL: Failed to connect to database: %v", err)
	}
	log.Println("Database connection verified")

	// Setup routes
	// Initialize productHandler even if db connection fails,
	// but ensure productRepo and productService handle nil db gracefully or panic early.
	// For now, we assume NewProductRepository can handle a nil db (e.g., by returning a no-op repo or erroring on method calls).
	// A more robust solution would be to not register these routes if db connection fails,
	// or to have a mock/fallback repository.
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	mux.HandleFunc("GET /api/products", productHandler.HandleProducts)
	mux.HandleFunc("POST /api/products", productHandler.HandleProducts)
	mux.HandleFunc("GET /api/products/", productHandler.HandleProductByID)
	mux.HandleFunc("PUT /api/products/", productHandler.HandleProductByID)
	mux.HandleFunc("DELETE /api/products/", productHandler.HandleProductByID)

	// Category Setup
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	mux.HandleFunc("GET /api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("POST /api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("GET /api/categories/", categoryHandler.HandleCategoryByID)
	mux.HandleFunc("PUT /api/categories/", categoryHandler.HandleCategoryByID)
	mux.HandleFunc("DELETE /api/categories/", categoryHandler.HandleCategoryByID)

	// Transaction routes
	mux.HandleFunc("POST /api/checkout", transactionHandler.HandleCheckout)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server up and running in", addr)

	serverErr := http.ListenAndServe(addr, handler)
	if serverErr != nil {
		fmt.Println("not running server", serverErr)
	}
}
