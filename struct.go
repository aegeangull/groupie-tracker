package main

import "time"

// A variable to map the entire artists response
var responseArtists []Artist

// A struct to map every artist to.
type Artist struct {
	Id              int               `json:"id"`
	Image           string            `json:"image"`
	Name            string            `json:"name"`
	Members         []string          `json:"members"`
	CreationDate    int               `json:"creationDate"`
	FirstAlbum      string            `json:"firstAlbum"`
	FirstAlbumDate  time.Time         `json:"firstAlbumDate"` // Will add it later
	LocationsAPI    string            `json:"locations"`
	ConcertDatesAPI string            `json:"concertDates"`
	RelationsAPI    string            `json:"relations"`
	Locations       []string          `json:"locations-slice"`    // Will add it later
	ConcertDates    []string          `json:"concertDates-slice"` // Will add it later
	DatesLocations  map[time.Time]string `json:"datesLocations"`  // Will add it from relations API later
}

type outInfo []Artist

type errorData struct {
	Num  int
	Text string
}

// A variable to map the entire relations response
var responseRelations Relations

// A struct to map every artists relation to.
type Relations struct {
	Index []Relation `json:"index"`
}

// A struct to map every relation "Location - Dates" to.
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"` // Map "location" string - "dates" slice of strings
}
