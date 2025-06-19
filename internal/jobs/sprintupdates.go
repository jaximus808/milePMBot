package jobs

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	jobs "github.com/jaximus808/milePMBot/internal/jobs/functions"
)

func StartSprintUpdateJob() (gocron.Scheduler, error) {

	location, err := time.LoadLocation("America/Chicago")

	if err != nil {
		return nil, err
	}

	s, err := gocron.NewScheduler(
		gocron.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(
		gocron.WeeklyJob(
			1,
			gocron.NewWeekdays(time.Tuesday),
			gocron.NewAtTimes(gocron.NewAtTime(16, 18, 0)),
		),
		gocron.NewTask(
			jobs.WeeklyRemindProjects,
		),
	)
	if err != nil {
		return nil, err
	}
	return s, nil

}
