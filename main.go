package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"kasir-api-irfan/database"
	_ "kasir-api-irfan/docs"
	"kasir-api-irfan/handlers"
	"kasir-api-irfan/repositories"
	"kasir-api-irfan/services"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Categories{
	{ID: 1, Name: "Elektronik", Description: "Perangkat gadget, komputer, dan aksesori elektronik lainnya"},
	{ID: 2, Name: "Pakaian", Description: "Berbagai jenis kain, baju, celana, dan aksesoris fashion"},
	{ID: 3, Name: "Makanan", Description: "Produk konsumsi siap saji maupun bahan makanan"},
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Mengambil semua data kategori
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/categories [get]
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "All categories retrieved successfully",
		Data:    categories,
	})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Mengambil kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Category retrieved successfully",
				Data:    c,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: "Category not found",
	})
}

// CreateCategory godoc
// @Summary Create new category
// @Description Menambahkan kategori baru
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Categories true "Category Data"
// @Success 201 {object} map[string]interface{}
// @Router /api/categories [post]
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newCategory Categories
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Category created successfully",
		Data:    categories,
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Menghapus kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/categories/{id} [delete]
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Category deleted successfully",
				Data:    categories,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: "Category not found",
	})
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update kategori berdasarkan ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Categories true "Category Data"
// @Success 200 {object} map[string]interface{}
// @Router /api/categories/{id} [put]
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)
	var updateCategory Categories
	err := json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories[i].Name = updateCategory.Name
			categories[i].Description = updateCategory.Description
			categories[i].ID = id
			categories = append(categories[:i], categories[i+1:]...)
			categories = append(categories, updateCategory)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Category updated successfully",
				Data:    categories,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: "Category not found",
	})
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

	//get all categories
	mux.HandleFunc("GET /api/categories", GetAllCategories)

	// create new category
	mux.HandleFunc("POST /api/categories", CreateCategory)

	// get category by id
	mux.HandleFunc("GET /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		GetCategoryByID(w, r)
	})

	//update category by id
	mux.HandleFunc("PUT /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		UpdateCategory(w, r)
	})

	//delete category by id
	mux.HandleFunc("DELETE /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		DeleteCategory(w, r)
	})

	// check api health
	mux.HandleFunc("/health", HealthCheck)

	// swagger documentation
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Apply CORS middleware to all routes
	handler := corsMiddleware(mux)

	// Initialize Database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
	}

	// Setup routes
	// Initialize productHandler even if db connection fails,
	// but ensure productRepo and productService handle nil db gracefully or panic early.
	// For now, we assume NewProductRepository can handle a nil db (e.g., by returning a no-op repo or erroring on method calls).
	// A more robust solution would be to not register these routes if db connection fails,
	// or to have a mock/fallback repository.
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Injeksi database routes ke mux (Selalu aktif agar tidak 404)
	mux.HandleFunc("GET /api/products", productHandler.HandleProducts)
	mux.HandleFunc("POST /api/products", productHandler.HandleProducts)
	mux.HandleFunc("GET /api/products/", productHandler.HandleProductByID)
	mux.HandleFunc("PUT /api/products/", productHandler.HandleProductByID)
	mux.HandleFunc("DELETE /api/products/", productHandler.HandleProductByID)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	serverErr := http.ListenAndServe(addr, handler)
	if serverErr != nil {
		fmt.Println("gagal running server", serverErr)
	}
}
