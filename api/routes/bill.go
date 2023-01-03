package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goRestAPI/model"
	"log"
	"strconv"
)

type Bill struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (b *Bill) getBill(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid bill ID"})
		return
	}

	bill := model.BillSchema{ID: id}
	if err = bill.QGetBill(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "bill not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, bill)
}

func (b *Bill) createBill(ctx *gin.Context, db *sql.DB) {
	var bill model.BillSchema
	if err := ctx.BindJSON(&bill); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := bill.QCreateBill(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new bill"})
		return
	}

	ctx.JSON(200, bill)
}

func (b *Bill) updateBill(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid bill id"})
		return
	}

	var bill model.BillSchema
	if err = ctx.BindJSON(&bill); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	bill.ID = id

	if err = bill.QUpdateBill(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "bill not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, bill)
}

func (b *Bill) deleteBill(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid bill id"})
		return
	}

	bill := model.BillSchema{ID: id}
	if err = bill.QDeleteBill(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "bill not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, bill)
}

func (b *Bill) getBills(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	var bill = model.BillSchema{}
	bills, err := bill.QGetBills(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, bills)
}

func (b *Bill) InitializeRoutes(db *sql.DB) {
	b.Router.GET("/bill", func(ctx *gin.Context) {
		b.getBills(ctx, db)
	})
	b.Router.POST("/bill", func(ctx *gin.Context) {
		b.createBill(ctx, db)
	})
	b.Router.GET("/bill/:id", func(ctx *gin.Context) {
		b.getBill(ctx, db)
	})
	b.Router.PUT("/bill/:id", func(ctx *gin.Context) {
		b.updateBill(ctx, db)
	})
	b.Router.DELETE("/bill/:id", func(ctx *gin.Context) {
		b.deleteBill(ctx, db)
	})
}
