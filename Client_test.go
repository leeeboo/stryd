package stryd

import (
	"testing"
	"time"

	"github.com/caarlos0/env"
)

func TestActivities(t *testing.T) {
	// TODO
	type Config struct {
		StrydEmail    string `env:"STRYD_EMAIL"`
		StrydPassword string `env:"STRYD_PASSWORD"`
	}

	cfg := Config{}
	err := env.Parse(&cfg)

	if err != nil {
		panic(err)
	}

	strydClient := NewClient(cfg.StrydEmail, cfg.StrydPassword)

	_, err = strydClient.Login()

	if err != nil {
		t.Fatal(err)
	}

	monthTime := getFirstDateOfMonth(time.Now())

	activities, err := strydClient.Activities(monthTime.Unix(), false)

	if err != nil {
		t.Fatal(err)
	}

	var distance float64

	for _, activity := range activities {
		distance += activity.Distance
	}

	distance = distance / 1000

	t.Logf("Distance: %.2f km", distance)
}

func getFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return getZeroTime(d)
}

func getZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
