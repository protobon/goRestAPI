package app

import (
	"awesomeProject/app/routes"
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

	var credit = routes.Credit{Router: a.Router, DB: a.DB}
	var product = routes.Product{Router: a.Router, DB: a.DB}
	credit.InitializeRoutes(credit.DB)
	product.InitializeRoutes(product.DB)
	schedule.RunCronJobs(a.DB)
}
