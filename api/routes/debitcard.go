package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goRestAPI/model"
	"log"
	"strconv"
)

type DebitCard struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (dc *DebitCard) getCard(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid DebitCard ID"})
		return
	}

	c := model.DebitCardSchema{ID: id}
	if err = c.QGetCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "DebitCard not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (dc *DebitCard) createCard(ctx *gin.Context, db *sql.DB) {
	var card model.DebitCardSchema
	if err := ctx.BindJSON(&card); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := card.QCreateCard(db); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Could not create new DebitCard"})
		return
	}

	ctx.JSON(200, card)
}

func (dc *DebitCard) updateCard(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid DebitCard id"})
		return
	}

	var card model.DebitCardSchema
	if err = ctx.BindJSON(&card); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	card.ID = id

	if err = card.QUpdateCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "DebitCard not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, card)
}

func (dc *DebitCard) deleteCard(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid DebitCard id"})
		return
	}

	card := model.DebitCardSchema{ID: id}
	if err = card.QDeleteCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Debit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, card)
}

func (dc *DebitCard) getCards(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var card = model.DebitCardSchema{}
	cards, err := card.QGetCards(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, cards)
}

func (dc *DebitCard) InitializeRoutes(db *sql.DB) {
	dc.Router.GET("/debit/cards", func(ctx *gin.Context) {
		dc.getCards(ctx, db)
	})
	dc.Router.POST("/debit/card", func(ctx *gin.Context) {
		dc.createCard(ctx, db)
	})
	dc.Router.GET("/debit/card/:id", func(ctx *gin.Context) {
		dc.getCard(ctx, db)
	})
	dc.Router.PUT("/debit/card/:id", func(ctx *gin.Context) {
		dc.updateCard(ctx, db)
	})
	dc.Router.DELETE("/debit/card/:id", func(ctx *gin.Context) {
		dc.deleteCard(ctx, db)
	})
}
