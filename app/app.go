package app

import (
	"awesomeProject/database"
	"awesomeProject/model"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *App) Run(addr string) {
	router := a.Router

	err := router.Run(addr)
	if err != nil {
		return
	}
}

func (a *App) Initialize(user string, password string, dbname string) {
	fmt.Println("Initializing App...")
	var err error
	a.DB = database.DBInit(user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = gin.Default()
	a.initializeRoutes()
}

func (a *App) getProduct(c *gin.Context) {
	productId := c.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(400, map[string]string{"error": "Invalid Product ID"})
		return
	}

	p := model.Product{ID: id}
	if err := p.GetProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(404, map[string]string{"error": "Product not found"})
		default:
			c.JSON(500, err.Error())
		}
		return
	}

	c.JSON(200, p)
}

func (a *App) createProduct(c *gin.Context) {
	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := p.CreateProduct(a.DB); err != nil {
		c.JSON(400, map[string]string{"error": "Could not create new Product"})
		return
	}

	c.JSON(200, p)
}

func (a *App) updateProduct(c *gin.Context) {
	productId := c.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(400, map[string]string{"error": "Invalid Product id"})
		return
	}

	var p model.Product
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	p.ID = id

	if err := p.UpdateProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(404, map[string]string{"error": "Product not found"})
		default:
			c.JSON(500, err.Error())
		}
		return
	}

	c.JSON(200, p)
}

func (a *App) deleteProduct(c *gin.Context) {
	productId := c.Param("id")
	id, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(400, map[string]string{"error": "Invalid Product id"})
		return
	}

	p := model.Product{ID: id}
	if err := p.DeleteProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(404, map[string]string{"error": "Product not found"})
		default:
			c.JSON(500, err.Error())
		}
		return
	}

	c.JSON(200, p)
}

func (a *App) getProducts(c *gin.Context) {
	start, _ := strconv.Atoi(c.Query("start"))
	count, _ := strconv.Atoi(c.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := model.GetProducts(a.DB, start, count)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, products)
}

func (a *App) initializeRoutes() {
	a.Router.GET("/products", func(c *gin.Context) {
		a.getProducts(c)
	})
	a.Router.POST("/product", func(c *gin.Context) {
		a.createProduct(c)
	})
	a.Router.GET("/product/:id", func(c *gin.Context) {
		a.getProduct(c)
	})
	a.Router.PUT("/product/:id", func(c *gin.Context) {
		a.updateProduct(c)
	})
	a.Router.DELETE("/product/:id", func(c *gin.Context) {
		a.deleteProduct(c)
	})
}
