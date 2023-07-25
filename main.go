package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/saaste/memento-mori/utils"
)

type AppConfig struct {
	Birthday       utils.Date
	LifeExpectancy int
	Events         []Event
}

type Event struct {
	Date  utils.Date
	Title string
	Label string
}

type Week struct {
	Class string
	Event Event
}

type Year struct {
	Year  int
	Weeks []Week
}

type IndexData struct {
	Birthday       utils.Date
	LifeExpectancy int
	Years          []Year
}

func getWeeks(year int, month int, config AppConfig) []Week {

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
		if startOfSquare.Before(config.Birthday.Time) {
			class = "hidden"
		} else if endOfSquare.Before(now) {
			class = "past"
		}

		newWeek := Week{
			Class: class,
			Event: Event{},
		}

		for _, event := range config.Events {
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

func getYears(yearOfBirth int, lifeExpectancy int, config AppConfig) []Year {
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	config, err := readConfig()
	if err != nil {
		fmt.Printf("ERROR: reading config.json failed: %s\n", err)
		os.Exit(1)
	}

	years := getYears(config.Birthday.Year(), config.LifeExpectancy, config)
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	data := IndexData{
		Birthday:       config.Birthday,
		LifeExpectancy: config.LifeExpectancy,
		Years:          years,
	}
	tmpl.Execute(w, data)
}

func readConfig() (AppConfig, error) {
	var config AppConfig

	f, err := os.ReadFile("./config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal([]byte(f), &config)
	if err != nil {
		return config, err
	}
	return config, nil

}

func main() {
	if !utils.FileExists("./config.json") {
		fmt.Println("ERROR: config.json does not exist")
		os.Exit(1)
	}

	if !utils.FileExists("./static/labels.css") {
		fmt.Println("ERROR: static/labels.css does not exist")
		os.Exit(1)
	}

	port := "3333"
	args := os.Args[1:]

	if len(args) > 0 {
		if _, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("ERROR: port must be an integer")
			os.Exit(1)
		}
		port = args[0]
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Server running on port %s\n", port)
	fmt.Printf("Go to http://localhost:%s to access your calendar\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("ERROR: unable to start the server: %s\n", err)
		os.Exit(1)
	}
}
