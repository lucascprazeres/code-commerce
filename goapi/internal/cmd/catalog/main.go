package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lucascprazeres/code-commerce/goapi/internal/database"
	"github.com/lucascprazeres/code-commerce/goapi/internal/service"
	"github.com/lucascprazeres/code-commerce/goapi/internal/webserver"
)

func main() {
	// database connection
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/code-commerce")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// setup services
	categoryDB := database.NewCategoryDB(db)
	categoryService := service.NewCategoryService(*categoryDB)

	productDB := database.NewProductDB(db)
	productService := service.NewProductService(*productDB)

	// setup handlers
	webCategoryHandler := webserver.NewWebCategoryHandler(categoryService)
	webProductHandler := webserver.NewProductHandler(productService)

	// setup router
	c := chi.NewRouter()

	c.Use(middleware.Logger)
	c.Use(middleware.Recoverer)

	c.Get("/category/{id}", webCategoryHandler.GetCategory)
	c.Get("/category", webCategoryHandler.GetCategories)
	c.Post("/category", webCategoryHandler.CreateCategory)

	c.Get("/product/{id}", webProductHandler.GetProduct)
	c.Get("/product", webProductHandler.GetProducts)
	c.Get("/product/category/{categoryID}", webProductHandler.GetProductByCategoryID)
	c.Post("/product", webProductHandler.CreateProduct)

	// run server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", c)
}
