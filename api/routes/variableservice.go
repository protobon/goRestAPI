package routes

import (
	"awesomeProject/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type ServiceVariable struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (service *ServiceVariable) getService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Variable Service ID"})
		return
	}

	s := model.ProductSchema{ID: id}
	if err = s.QGetProduct(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Variable Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceVariable) insertService(ctx *gin.Context, db *sql.DB) {
	var s model.FixedServiceSchema
	if err := ctx.BindJSON(&s); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := s.QCreateService(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Variable Service"})
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceVariable) updateService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Variable Service id"})
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
			ctx.JSON(404, map[string]string{"error": "Variable Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceVariable) deleteService(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Variable Service id"})
		return
	}

	s := model.FixedServiceSchema{ID: id}
	if err = s.QDeleteService(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Variable Service not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, s)
}

func (service *ServiceVariable) getServices(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := model.QGetVariableServices(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, products)
}

func (service *ServiceVariable) InitializeRoutes(db *sql.DB) {
	service.Router.GET("/services/variable", func(ctx *gin.Context) {
		service.getServices(ctx, db)
	})
	service.Router.POST("/service/variable", func(ctx *gin.Context) {
		service.insertService(ctx, db)
	})
	service.Router.GET("/service/variable/:id", func(ctx *gin.Context) {
		service.getService(ctx, db)
	})
	service.Router.PUT("/service/variable/:id", func(ctx *gin.Context) {
		service.updateService(ctx, db)
	})
	service.Router.DELETE("/service/variable/:id", func(ctx *gin.Context) {
		service.deleteService(ctx, db)
	})
}
