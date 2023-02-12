package routes

import (
	"database/sql"
	"goRestAPI/model"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Dummy struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (d *Dummy) getDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy ID"})
		return
	}

	dummy := model.DummySchema{ID: id}
	if err = dummy.QGetDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *Dummy) createDummy(ctx *gin.Context, db *sql.DB) {
	var dummy model.DummySchema
	if err := ctx.BindJSON(&dummy); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := dummy.QCreateDummy(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Dummy"})
		return
	}

	ctx.JSON(200, dummy)
}

func (d *Dummy) updateDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy id"})
		return
	}

	var dummy model.DummySchema
	if err = ctx.BindJSON(&dummy); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	dummy.ID = id

	if err = dummy.QUpdateDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *Dummy) deleteDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy id"})
		return
	}

	dummy := model.DummySchema{ID: id}
	if err = dummy.QDeleteDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *Dummy) getDummies(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var dummy = model.DummySchema{}
	dummies, err := dummy.QGetDummies(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, dummies)
}

func (d *Dummy) InitializeRoutes(db *sql.DB) {
	d.Router.GET("/dummies", func(ctx *gin.Context) {
		d.getDummies(ctx, db)
	})
	d.Router.POST("/dummy", func(ctx *gin.Context) {
		d.createDummy(ctx, db)
	})
	d.Router.GET("/dummy/:id", func(ctx *gin.Context) {
		d.getDummy(ctx, db)
	})
	d.Router.PUT("/dummy/:id", func(ctx *gin.Context) {
		d.updateDummy(ctx, db)
	})
	d.Router.DELETE("/dummy/:id", func(ctx *gin.Context) {
		d.deleteDummy(ctx, db)
	})
}
