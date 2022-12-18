package api

import (
	"awesomeProject/api/routes"
	"awesomeProject/database"
	"awesomeProject/schedule"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *App) Run(addr string) {
	router := a.Router

	err := router.Run(addr)
	if err != nil {
		return
	}
}

func (a *App) Initialize(user string, password string, dbname string) {
	fmt.Println("Initializing App...")
	var err error
	a.DB = database.DBInit(user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = gin.Default()

	var product = routes.Product{Router: a.Router, DB: a.DB}
	var card = routes.Card{Router: a.Router, DB: a.DB}
	var credit = routes.Credit{Router: a.Router, DB: a.DB}
	var serviceFixed = routes.ServiceFixed{Router: a.Router, DB: a.DB}
	var serviceVariable = routes.ServiceVariable{Router: a.Router, DB: a.DB}

	card.InitializeRoutes(card.DB)
	product.InitializeRoutes(product.DB)
	credit.InitializeRoutes(credit.DB)
	serviceFixed.InitializeRoutes(serviceFixed.DB)
	serviceVariable.InitializeRoutes(serviceVariable.DB)
	schedule.RunCronJobs(a.DB)
}
