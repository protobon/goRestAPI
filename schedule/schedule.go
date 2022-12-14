package schedule

import (
	"awesomeProject/model"
	"database/sql"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

func RunCronJobs(db *sql.DB) {
	s := gocron.NewScheduler(time.UTC)
	fmt.Println("Schedule QNextQuota() Every 1st of each month...")
	_, err := s.Every(1).Month(1).Do(func() {
		c := model.CreditSchema{}
		err := c.QNextQuota(db)
		if err != nil {
			return
		}
	})
	if err != nil {
		return
	}
	fmt.Println("Scheduled.")
}
