package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goRestAPI/model"
	"log"
	"strconv"
)

type Rent struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (r *Rent) getRent(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Rent ID"})
		return
	}

	rent := model.RentSchema{ID: id}
	if err = rent.QGetRent(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Rent not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, rent)
}

func (r *Rent) insertRent(ctx *gin.Context, db *sql.DB) {
	var rent model.RentSchema
	if err := ctx.BindJSON(&rent); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := rent.QCreateRent(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Rent"})
		return
	}

	ctx.JSON(200, rent)
}

func (r *Rent) updateRent(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Rent id"})
		return
	}

	var rent model.RentSchema
	if err = ctx.BindJSON(&rent); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	rent.ID = id

	if err = rent.QUpdateRent(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Rent not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, rent)
}

func (r *Rent) deleteRent(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Rent id"})
		return
	}

	rent := model.RentSchema{ID: id}
	if err = rent.QDeleteRent(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Rent not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, rent)
}

func (r *Rent) InitializeRoutes(db *sql.DB) {
	r.Router.POST("/rent", func(ctx *gin.Context) {
		r.insertRent(ctx, db)
	})
	r.Router.GET("/rent/:id", func(ctx *gin.Context) {
		r.getRent(ctx, db)
	})
	r.Router.PUT("/rent/:id", func(ctx *gin.Context) {
		r.updateRent(ctx, db)
	})
	r.Router.DELETE("/rent/:id", func(ctx *gin.Context) {
		r.deleteRent(ctx, db)
	})
}
