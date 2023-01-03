package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goRestAPI/model"
	"log"
	"strconv"
)

type Debit struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (d *Debit) getDebit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Debit ID"})
		return
	}

	debit := model.DebitSchema{ID: id}
	if err = debit.QGetDebit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Debit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, debit)
}

func (d *Debit) createDebit(ctx *gin.Context, db *sql.DB) {
	var debit model.DebitSchema
	if err := ctx.BindJSON(&debit); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := debit.QCreateDebit(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Debit"})
		return
	}

	ctx.JSON(200, debit)
}

func (d *Debit) deleteDebit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Debit id"})
		return
	}

	debit := model.DebitSchema{ID: id}
	if err = debit.QDeleteDebit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Debit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, d)
}

func (d *Debit) getDebits(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var debit = model.DebitSchema{}
	debits, err := debit.QGetDebits(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, debits)
}

func (d *Debit) InitializeRoutes(db *sql.DB) {
	d.Router.GET("/debits", func(ctx *gin.Context) {
		d.getDebits(ctx, db)
	})
	d.Router.POST("/debit", func(ctx *gin.Context) {
		d.createDebit(ctx, db)
	})
	d.Router.GET("/debit/:id", func(ctx *gin.Context) {
		d.getDebit(ctx, db)
	})
	d.Router.DELETE("/debit/:id", func(ctx *gin.Context) {
		d.deleteDebit(ctx, db)
	})
}
