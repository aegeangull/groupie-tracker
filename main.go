package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var defaultInfo Artist

func main() {
	loadArtists()
	loadConcerts()
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/artist_card", artistCard)
	http.HandleFunc("/search", search)
	http.HandleFunc("/filter", filter)

	fmt.Println("Open http://localhost:8081")
	fmt.Println("Press ctrl+C to exit")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	tmpl.Execute(w, responseArtists)
}

func artistCard(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("template/pages/artist_card.html")
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	name := r.FormValue("name")
	for _, v := range responseArtists {
		if v.Name == name {
			tmpl.Execute(w, v)
			break
		}
	}
	return
}

func search(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("searchBar")
	tmpl, err := template.ParseFiles("template/pages/search.html")
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	var out outInfo
	var result []int
	for i := 0; i < 52; i++ {
		if strings.ReplaceAll((strings.ToLower(value)), " - band/artist", "") == strings.ToLower(responseArtists[i].Name) {
			result = append(result, i)
		}
		if strings.ReplaceAll((strings.ToLower(value)), " - first album date", "") == strings.ToLower(responseArtists[i].FirstAlbum) {
			result = append(result, i)
		}
		if strings.TrimRight(value, " - creation date") == strconv.Itoa(responseArtists[i].CreationDate) {
			result = append(result, i)
		}
		for _, v := range responseArtists[i].Members {
			if strings.ReplaceAll(strings.ToLower(value), " - member", "") == strings.ToLower(v) {
				result = append(result, i)
			}
		}
		for _, k := range responseArtists[i].Locations {
			if strings.ReplaceAll(strings.ToLower(value), " - concert location", "") == strings.ToLower(k) {
				result = append(result, i)
			}
		}
	}
	if len(result) < 1 {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	result = removeDuplicateInt(result)

	for _, v := range result {
		out = append(out, responseArtists[v])
	}
	tmpl.Execute(w, out)
	return
}

func filter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Required if you don't call r.FormValue()

	fromYearString := r.FormValue("fromYear")
	toYearString := r.FormValue("toYear")
	fromAlbumString := r.FormValue("fromAlbum")
	toAlbumString := r.FormValue("toAlbum")
	location := r.FormValue("loc")
	nMembersStrings := r.Form["members"]

	fromYear, err := strconv.Atoi(fromYearString)	// convert string to int
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	toYear, err := strconv.Atoi(toYearString)	// convert string to int
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	fromAlbum, err := strconv.Atoi(fromAlbumString)	// convert string to int
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	toAlbum, err := strconv.Atoi(toAlbumString)	// convert string to int
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var nMembers []int

	for i := 0; i < len(nMembersStrings); i++ {	// convert string to int
		n, err := strconv.Atoi(nMembersStrings[i])
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		nMembers = append(nMembers, n)
	}
	
	tmpl, err := template.ParseFiles("template/pages/filter.html")
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	var out outInfo
	var result []int
	
	for i := 0; i < 52; i++ {
		if responseArtists[i].CreationDate < fromYear || responseArtists[i].CreationDate > toYear {	// if creation date is out of range
			continue
		}
		albumYear, err := strconv.Atoi(responseArtists[i].FirstAlbum[6:])	// convert first album year string to int
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		if albumYear < fromAlbum || albumYear > toAlbum {	// if first album year is out of range
			continue
		}
		isLoc := true	// concert location is match by default
		if location != "-" {	// if we filter concert location
			isLoc = false
			for _, l := range responseArtists[i].DatesLocations {
				if l == location {
					isLoc = true
					break
				}
			}
		}
		if !isLoc {
			continue
		}
		mem := true	// number of members is match by default
		if len(nMembers) > 0 {	// if we filter number of memebers
			mem = false
			for n := 0; n < len(nMembers); n++ {
				if len(responseArtists[i].Members) == nMembers[n] {
					mem = true
				}
			}
		}
		if mem {
			result = append(result, i)
		}
	}
	if len(result) < 1 {
		tmpl.Execute(w, out)
		return
	}

	result = removeDuplicateInt(result)

	for _, v := range result {
		out = append(out, responseArtists[v])
	}
	tmpl.Execute(w, out)
	return
}

func errorHandler(w http.ResponseWriter, r *http.Request, err_num int) {
	tmpl, err := template.ParseFiles("template/pages/error.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	w.WriteHeader(err_num)
	errData := errorData{Num: err_num}
	if err_num == 404 {
		errData.Text = "Page Not Found"
	} else if err_num == 400 {
		errData.Text = "Bad Request"
	} else if err_num == 500 {
		errData.Text = "Internal Server Error"
	}
	tmpl.Execute(w, errData)
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func loadArtists() { // Consume artists API
	response := Request("https://groupietrackers.herokuapp.com/api/artists")
	if err := json.Unmarshal([]byte(response), &responseArtists); err != nil {
		panic(err)
	}
	for i := range responseArtists { // Add reformated FirstAlbumDate to Artist
		dateTime, _ := time.Parse("02-01-2006", responseArtists[i].FirstAlbum)
		responseArtists[i].FirstAlbumDate = dateTime
	}
}

func loadConcerts() { // Consume relation API
	response := Request("https://groupietrackers.herokuapp.com/api/relation")
	if err := json.Unmarshal([]byte(response), &responseRelations); err != nil {
		panic(err)
	}
	addConcertsToArtists() // Reformat relations and add to Artist struct
}

func Request(API string) []byte {
	response, err := http.Get(API) // query API endpoint
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(response.Body) // read []byte response
	if err != nil {
		fmt.Println(err)
	}

	defer func(Body io.ReadCloser) { // Close reader
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	return body
}

// Reformat relations and add to Artist JSON structure
func addConcertsToArtists() {
	for i := range responseRelations.Index { // For each Artist
		concerts := make(map[time.Time]string) // Create  "concerts" map, "date" Time - "location" string
		concertDates := []string{}             // Create concertDates string slice
		locations := []string{}                // Create locations string slice

		for loc, dates := range responseRelations.Index[i].DatesLocations { // For each location
			locations = append(locations, loc) // Add location to slice
			location := refLocation(loc)       // Reformat location

			for _, date := range dates { // For each date in location
				concertDates = append(concertDates, date)     // Add date to slice
				dateTime, _ := time.Parse("02-01-2006", date) // Reformat date
				concerts[dateTime] = location                 // Add "date-location" element to "concerts" map
			}
		}
		responseArtists[i].Locations = locations       // Add locations to Artist
		responseArtists[i].ConcertDates = concertDates // Add concertDates to Artist
		responseArtists[i].DatesLocations = concerts   // Add "concerts" map to Artist
	}
}

func refLocation(s string) string { // Reformat location
	return strings.ReplaceAll(strings.ToUpper(strings.ReplaceAll(s, "_", " ")), "-", ", ")
}
