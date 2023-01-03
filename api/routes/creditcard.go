package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goRestAPI/model"
	"log"
	"strconv"
)

type CreditCard struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (cc *CreditCard) getCard(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid CreditCard ID"})
		return
	}

	card := model.CreditCardSchema{ID: id}
	if err = card.QGetCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "CreditCard not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, cc)
}

func (cc *CreditCard) createCard(ctx *gin.Context, db *sql.DB) {
	var card model.CreditCardSchema
	if err := ctx.BindJSON(&card); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := card.CreateCard(db); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Could not create new CreditCard"})
		return
	}

	ctx.JSON(200, card)
}

func (cc *CreditCard) updateLimit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid CreditCard id"})
		return
	}

	var card model.CreditCardSchema
	if err := ctx.BindJSON(&card); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	card.ID = id

	if err = card.QUpdateLimit(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "CreditCard not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, card)
}

func (cc *CreditCard) deleteCard(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid CreditCard id"})
		return
	}

	card := model.CreditCardSchema{ID: id}
	if err = card.DeleteCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, card)
}

func (cc *CreditCard) getCards(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var card model.CreditCardSchema
	cards, err := card.GetCards(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, cards)
}

func (cc *CreditCard) InitializeRoutes(db *sql.DB) {
	cc.Router.GET("/credit/cards", func(ctx *gin.Context) {
		cc.getCards(ctx, db)
	})
	cc.Router.POST("/credit/card", func(ctx *gin.Context) {
		cc.createCard(ctx, db)
	})
	cc.Router.GET("/credit/card/:id", func(ctx *gin.Context) {
		cc.getCard(ctx, db)
	})
	cc.Router.PUT("/credit/card/:id", func(ctx *gin.Context) {
		cc.updateLimit(ctx, db)
	})
	cc.Router.DELETE("/credit/card/:id", func(ctx *gin.Context) {
		cc.deleteCard(ctx, db)
	})
}
