package routes

import (
	"awesomeProject/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type Card struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (card *Card) getCard(ctx *gin.Context, db *sql.DB) {
	cardId := ctx.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Card ID"})
		return
	}

	c := model.CardSchema{ID: id}
	if err := c.GetCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Card not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (card *Card) createCard(ctx *gin.Context, db *sql.DB) {
	var c model.CardSchema
	if err := ctx.BindJSON(&c); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := c.CreateCard(db); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Could not create new Card"})
		return
	}

	ctx.JSON(200, c)
}

func (card *Card) updateCard(ctx *gin.Context, db *sql.DB) {
	cardId := ctx.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Card id"})
		return
	}

	var c model.CardSchema
	if err := ctx.BindJSON(&c); err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	c.ID = id

	if err := c.UpdateCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Card not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (card *Card) deleteCard(ctx *gin.Context, db *sql.DB) {
	cardId := ctx.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Card id"})
		return
	}

	c := model.CardSchema{ID: id}
	if err := c.DeleteCard(db); err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Product not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (card *Card) getCards(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	cards, err := model.GetCards(db, start, count)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, cards)
}

func (card *Card) InitializeRoutes(db *sql.DB) {
	card.Router.GET("/cards", func(ctx *gin.Context) {
		card.getCards(ctx, db)
	})
	card.Router.POST("/card", func(ctx *gin.Context) {
		card.createCard(ctx, db)
	})
	card.Router.GET("/card/:id", func(ctx *gin.Context) {
		card.getCard(ctx, db)
	})
	card.Router.PUT("/card/:id", func(ctx *gin.Context) {
		card.updateCard(ctx, db)
	})
	card.Router.DELETE("/card/:id", func(ctx *gin.Context) {
		card.deleteCard(ctx, db)
	})
}
