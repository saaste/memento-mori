package calendar

import "github.com/saaste/memento-mori/config"

type Week struct {
	Class string
	Event config.Event
}

type Year struct {
	Year  int
	Weeks []Week
}
