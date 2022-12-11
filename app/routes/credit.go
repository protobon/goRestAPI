package routes

import (
	"awesomeProject/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Credit struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (credit *Credit) getCredit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Credit ID"})
		return
	}

	c := model.CreditSchema{ID: uint32(id)}
	if err := c.QGetCredit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Credit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (credit *Credit) createCredit(ctx *gin.Context, db *sql.DB) {
	var c model.CreditSchema
	if err := ctx.BindJSON(&c); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := c.QCreateCredit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Could not create new Credit"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (credit *Credit) updateCredit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Credit ID"})
		return
	}

	c := model.CreditSchema{ID: uint32(id)}
	if err := c.QUpdateCredit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Credit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (credit *Credit) deleteCredit(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Credit ID"})
		return
	}

	c := model.CreditSchema{ID: uint32(id)}
	if err := c.QDeleteCredit(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Credit not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, c)
}

func (credit *Credit) getCredits(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	credits, err := model.QGetCredits(db, start, count)

	if err != nil {
		ctx.JSON(500, err.Error())
	}

	ctx.JSON(200, credits)
}

func (credit *Credit) clearCredit(ctx *gin.Context, db *sql.DB) {
	if err := model.QClearCredit(db); err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, map[string]string{"message": "Credit table cleared"})
}

func (credit *Credit) InitializeRoutes(db *sql.DB) {
	credit.Router.GET("/credits", func(ctx *gin.Context) {
		credit.getCredits(ctx, db)
	})
	credit.Router.POST("/credit", func(ctx *gin.Context) {
		credit.createCredit(ctx, db)
	})
	credit.Router.GET("/credit/:id", func(ctx *gin.Context) {
		credit.getCredit(ctx, db)
	})
	credit.Router.PUT("/credit/:id", func(ctx *gin.Context) {
		credit.updateCredit(ctx, db)
	})
	credit.Router.DELETE("/credit/:id", func(ctx *gin.Context) {
		credit.deleteCredit(ctx, db)
	})
	//credit.Router.DELETE("/credit/clear", func(ctx *gin.Context) {
	//	credit.clearCredit(ctx, db)
	//})
}
