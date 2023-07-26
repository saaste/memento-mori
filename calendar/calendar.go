package calendar

import (
	"time"

	"github.com/saaste/memento-mori/config"
	"github.com/saaste/memento-mori/utils"
)

func getWeeks(year int, month int, conf config.AppConfig) []Week {

	now := time.Now()
	daysInMonth := utils.DaysInMonth(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC))
	daysInSquare := daysInMonth / 4

	weeks := make([]Week, 0)
	for w := 0; w < 4; w++ {
		startOfSquare := time.Date(year, time.Month(month), daysInSquare*w+1, 0, 0, 0, 0, time.UTC)
		endOfSquare := time.Date(year, time.Month(month), daysInSquare*(w+1), 0, 0, 0, 0, time.UTC)
		if w == 3 && (endOfSquare.Day() != daysInMonth || endOfSquare.Month() != startOfSquare.Month()) {
			endOfSquare = time.Date(year, time.Month(month), daysInMonth, 0, 0, 0, 0, time.UTC)
		}

		class := ""
		if startOfSquare.Before(conf.Birthday.Time) {
			class = "hidden"
		} else if endOfSquare.Before(now) {
			class = "past"
		}

		newWeek := Week{
			Class: class,
			Event: config.Event{},
		}

		for _, event := range conf.Events {
			if event.Date.Equal(startOfSquare) || event.Date.Equal(endOfSquare) || (event.Date.After(startOfSquare) && event.Date.Before(endOfSquare)) {
				year = event.Date.Year()
				newWeek.Event = event
				break
			}
		}

		weeks = append(weeks, newWeek)
	}
	return weeks
}

func GetYears(yearOfBirth int, lifeExpectancy int, config config.AppConfig) []Year {
	firstYear := yearOfBirth
	lastYear := yearOfBirth + lifeExpectancy

	years := make([]Year, 0)
	for currentYear := firstYear; currentYear <= lastYear; currentYear++ {
		weeks := make([]Week, 0)
		for currentMonth := 1; currentMonth < 12; currentMonth++ {
			weeks = append(weeks, getWeeks(currentYear, currentMonth, config)...)
		}
		years = append(years, Year{
			Year:  currentYear,
			Weeks: weeks,
		})
	}
	return years
}
