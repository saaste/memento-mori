package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/saaste/memento-mori/calendar"
	"github.com/saaste/memento-mori/config"
	"github.com/saaste/memento-mori/utils"
)

type IndexData struct {
	Birthday       utils.Date
	LifeExpectancy int
	Years          []calendar.Year
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("ERROR: reading config.json failed: %s\n", err)
	}

	years := calendar.GetYears(config.Birthday.Year(), config.LifeExpectancy, config)
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	data := IndexData{
		Birthday:       config.Birthday,
		LifeExpectancy: config.LifeExpectancy,
		Years:          years,
	}
	tmpl.Execute(w, data)
}

func main() {
	if !utils.FileExists("./config.json") {
		log.Fatalln("ERROR: config.json does not exist")
	}

	if !utils.FileExists("./static/labels.css") {
		log.Fatalln("ERROR: static/labels.css does not exist")
	}

	port := "3333"
	args := os.Args[1:]

	if len(args) > 0 {
		if _, err := strconv.Atoi(args[0]); err != nil {
			log.Fatalln("ERROR: port must be an integer")
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
		log.Fatalf("ERROR: unable to start the server: %s\n", err)
	}
}
