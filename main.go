package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api-irfan/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Product struct {
	ID    int    `json:"id:"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 1000, Stock: 10},
	{ID: 2, Name: "Smartphone", Price: 500, Stock: 20},
	{ID: 3, Name: "Tablet", Price: 300, Stock: 15},
}

var categories = []Categories{
	{ID: 1, Name: "Elektronik", Description: "Perangkat gadget, komputer, dan aksesori elektronik lainnya"},
	{ID: 2, Name: "Pakaian", Description: "Berbagai jenis kain, baju, celana, dan aksesoris fashion"},
	{ID: 3, Name: "Makanan", Description: "Produk konsumsi siap saji maupun bahan makanan"},
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

// GetAllProducts godoc
// @Summary Get all products
// @Description Mengambil semua data produk
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/products [get]
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "All products retrieved successfully",
		Data:    products,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Mengambil produk berdasarkan ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Product retrieved successfully",
				Data:    product,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: "Product not found",
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

// CreateProduct godoc
// @Summary Create new product
// @Description Menambahkan produk baru
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product Data"
// @Success 201 {object} map[string]interface{}
// @Router /api/products [post]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "product created successfully",
		Data:    products,
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

// UpdateProduct godoc
// @Summary Update product
// @Description Update produk berdasarkan ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Product Data"
// @Success 200 {object} map[string]interface{}
// @Router /api/products/{id} [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)
	var UpdateProduct Product
	err := json.NewDecoder(r.Body).Decode(&UpdateProduct)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			// Update product data
			products[i].Name = UpdateProduct.Name
			products[i].Price = UpdateProduct.Price
			products[i].Stock = UpdateProduct.Stock

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Product updated successfully",
				Data:    products[i],
			})
			return
		}
	}

}

// DeleteProduct godoc
// @Summary Delete product
// @Description Menghapus produk berdasarkan ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [delete]
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			// Remove product from slice
			products = append(products[:i], products[i+1:]...)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Product deleted successfully",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: "Product not found",
	})
}

// @title Kasir API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {

	//get all categories
	http.HandleFunc("GET /api/categories", GetAllCategories)

	// get all products
	http.HandleFunc("GET /api/products", GetAllProducts)

	// create new category
	http.HandleFunc("POST /api/categories", CreateCategory)

	// create new product
	http.HandleFunc("POST /api/products", CreateProduct)

	// get category by id
	http.HandleFunc("GET /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		GetCategoryByID(w, r)
	})

	//get all products by id
	http.HandleFunc("GET /api/products/", func(w http.ResponseWriter, r *http.Request) {
		GetProductByID(w, r)
	})

	//update category by id
	http.HandleFunc("PUT /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		UpdateCategory(w, r)
	})

	//update product by id
	http.HandleFunc("PUT /api/products/", func(w http.ResponseWriter, r *http.Request) {
		UpdateProduct(w, r)
	})

	//delete product by id
	http.HandleFunc("DELETE /api/products/", func(w http.ResponseWriter, r *http.Request) {
		DeleteProduct(w, r)
	})

	//delete category by id
	http.HandleFunc("DELETE /api/categories/", func(w http.ResponseWriter, r *http.Request) {
		DeleteCategory(w, r)
	})

	// check api health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status:  "success",
			Message: "API is healthy",
		})
	})

	// swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server is running...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
