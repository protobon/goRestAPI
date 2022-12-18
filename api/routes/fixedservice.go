package routes

import (
	"awesomeProject/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type ServiceFixed struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (service *ServiceFixed) getService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Fixed Service ID"})
		return
	}

	s := model.ProductSchema{ID: id}
	if err = s.QGetProduct(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Fixed Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceFixed) insertService(ctx *gin.Context, db *sql.DB) {
	var s model.FixedServiceSchema
	if err := ctx.BindJSON(&s); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := s.QCreateService(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Fixed Service"})
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceFixed) updateService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Fixed Service id"})
		return
	}

	var s model.FixedServiceSchema
	if err = ctx.BindJSON(&s); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	s.ID = id

	if err = s.QUpdateService(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Fixed Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceFixed) deleteService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Fixed Service id"})
		return
	}

	s := model.FixedServiceSchema{ID: id}
	if err = s.QDeleteService(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Fixed Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceFixed) getServices(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := model.QGetFixedServices(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, products)
}

func (service *ServiceFixed) InitializeRoutes(db *sql.DB) {
	service.Router.GET("/services/fixed", func(ctx *gin.Context) {
		service.getServices(ctx, db)
	})
	service.Router.POST("/service/fixed", func(ctx *gin.Context) {
		service.insertService(ctx, db)
	})
	service.Router.GET("/service/fixed/:id", func(ctx *gin.Context) {
		service.getService(ctx, db)
	})
	service.Router.PUT("/service/fixed/:id", func(ctx *gin.Context) {
		service.updateService(ctx, db)
	})
	service.Router.DELETE("/service/fixed/:id", func(ctx *gin.Context) {
		service.deleteService(ctx, db)
	})
}
