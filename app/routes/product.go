package routes

import (
	"awesomeProject/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type Product struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (product *Product) getProduct(ctx *gin.Context, db *sql.DB) {
	productId := ctx.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Product ID"})
		return
	}

	p := model.ProductSchema{ID: id}
	if err := p.GetProduct(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Product not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, p)
}

func (product *Product) createProduct(ctx *gin.Context, db *sql.DB) {
	var p model.ProductSchema
	if err := ctx.BindJSON(&p); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := p.CreateProduct(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Product"})
		return
	}

	ctx.JSON(200, p)
}

func (product *Product) updateProduct(ctx *gin.Context, db *sql.DB) {
	productId := ctx.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Product id"})
		return
	}

	var p model.ProductSchema
	if err := ctx.BindJSON(&p); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	p.ID = id

	if err := p.UpdateProduct(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Product not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, p)
}

func (product *Product) deleteProduct(ctx *gin.Context, db *sql.DB) {
	productId := ctx.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Product id"})
		return
	}

	p := model.ProductSchema{ID: id}
	if err := p.DeleteProduct(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Product not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, p)
}

func (product *Product) getProducts(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := model.GetProducts(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, products)
}

func (product *Product) InitializeRoutes(db *sql.DB) {
	product.Router.GET("/products", func(ctx *gin.Context) {
		product.getProducts(ctx, db)
	})
	product.Router.POST("/product", func(ctx *gin.Context) {
		product.createProduct(ctx, db)
	})
	product.Router.GET("/product/:id", func(ctx *gin.Context) {
		product.getProduct(ctx, db)
	})
	product.Router.PUT("/product/:id", func(ctx *gin.Context) {
		product.updateProduct(ctx, db)
	})
	product.Router.DELETE("/product/:id", func(ctx *gin.Context) {
		product.deleteProduct(ctx, db)
	})
}
