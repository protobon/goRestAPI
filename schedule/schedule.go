package schedule

import (
	"database/sql"
	"fmt"
	"github.com/go-co-op/gocron"
	"goRestAPI/model"
	"time"
)

func RunCronJobs(db *sql.DB) {
	s := gocron.NewScheduler(time.UTC)
	fmt.Println("Schedule QNextQuota() Every 1st of each month...")
	_, err := s.Every(1).Month(1).Do(func() {
		err := model.QCreditNextQuotaAll(db)
		if err != nil {
			return
		}
	})
	if err != nil {
		return
	}
	fmt.Println("Scheduled.")
}
