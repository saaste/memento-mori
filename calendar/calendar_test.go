package calendar

import (
	"testing"
	"time"

	"github.com/saaste/memento-mori/config"
	"github.com/saaste/memento-mori/utils"
)

func TestGetWeeksWithEvents(t *testing.T) {
	events := []config.Event{
		{
			Date:  utils.Date{Time: time.Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC)},
			Title: "First",
		},
		{
			Date:  utils.Date{Time: time.Date(2023, time.January, 14, 0, 0, 0, 0, time.UTC)},
			Title: "Second",
		},
		{
			Date:  utils.Date{Time: time.Date(2023, time.January, 21, 0, 0, 0, 0, time.UTC)},
			Title: "Third",
		},
		{
			Date:  utils.Date{Time: time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC)},
			Title: "Forth",
		},
	}

	conf := config.AppConfig{
		Birthday:       utils.Date{Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
		LifeExpectancy: 50,
		Events:         events,
	}

	weeks := getWeeks(2023, 1, conf)
	if len(weeks) != 4 {
		t.Fatalf("Expected month to have 4 weeks, actual: %d", len(weeks))
	}

	if weeks[0].Event.Title != events[0].Title {
		t.Fatalf("Expected first event title to be %s, actual: %s", events[0].Title, weeks[0].Event.Title)
	}

	if weeks[1].Event.Title != events[1].Title {
		t.Fatalf("Expected second event title to be %s, actual: %s", events[1].Title, weeks[1].Event.Title)
	}

	if weeks[2].Event.Title != events[2].Title {
		t.Fatalf("Expected third event title to be %s, actual: %s", events[2].Title, weeks[2].Event.Title)
	}

	if weeks[3].Event.Title != events[3].Title {
		t.Fatalf("Expected fourth event title to be %s, actual: %s", events[3].Title, weeks[3].Event.Title)
	}
}

func TestGetWeeksWithHiddenWeeks(t *testing.T) {
	conf := config.AppConfig{
		Birthday:       utils.Date{Time: time.Date(2000, time.January, 14, 0, 0, 0, 0, time.UTC)},
		LifeExpectancy: 50,
		Events:         make([]config.Event, 0),
	}
	weeks := getWeeks(2000, 1, conf)

	if weeks[0].Class != "hidden" {
		t.Fatalf("Expected first week class to be 'past', actual: %s", weeks[0].Class)
	}

	if weeks[1].Class != "hidden" {
		t.Fatalf("Expected second week class to be 'past', actual: %s", weeks[1].Class)
	}

	if weeks[2].Class == "hidden" {
		t.Fatalf("Expected third week class to not be 'hidden', actual: %s", weeks[2].Class)
	}

	if weeks[3].Class == "hidden" {
		t.Fatalf("Expected fourth week class to not be 'hidden', actual: %s", weeks[3].Class)
	}
}

func TestYears(t *testing.T) {
	conf := config.AppConfig{
		Birthday:       utils.Date{Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
		LifeExpectancy: 5,
		Events:         make([]config.Event, 0),
	}
	years := GetYears(2000, 5, conf)

	if len(years) != 6 {
		t.Fatalf("Expected 5 years, actual: %d", len(years))
	}

	for i, year := range years {
		if year.Year != conf.Birthday.Year()+i {
			t.Fatalf("Expected year[%d].Year to be %d, actual: %d", i, conf.Birthday.Year()+i, year.Year)
		}
		if len(year.Weeks) != 48 {
			t.Fatalf("Expected year[%d] to have 48 weeks, actual: %d", i, len(year.Weeks))
		}
	}
}
