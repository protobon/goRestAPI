package api

import (
	"database/sql"
	"fmt"
	"goRestAPI/api/routes"
	"goRestAPI/database"
	"log"

	"github.com/gin-gonic/gin"
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

	var dummy = routes.Dummy{Router: a.Router, DB: a.DB}
	// var creditCard = routes.CreditCard{Router: a.Router, DB: a.DB}
	// var credit = routes.Credit{Router: a.Router, DB: a.DB}
	// var rent = routes.Rent{Router: a.Router, DB: a.DB}
	// var bill = routes.Bill{Router: a.Router, DB: a.DB}
	// var debitCard = routes.DebitCard{Router: a.Router, DB: a.DB}

	dummy.InitializeRoutes(dummy.DB)
	// creditCard.InitializeRoutes(creditCard.DB)
	// credit.InitializeRoutes(credit.DB)
	// rent.InitializeRoutes(rent.DB)
	// bill.InitializeRoutes(bill.DB)
	// debitCard.InitializeRoutes(debitCard.DB)
	// schedule.RunCronJobs(a.DB)
	fmt.Println("***** App Running *****")
}
