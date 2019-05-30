package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// App struct
type App struct {
	URLTemplate string
	Accept      string
	Query       string
	Source      string
	Size        string
	Offset      string
	Client      http.Client
}

// Body struct
type Body struct {
	Result Result `json:"result"`
}

// Result struct
type Result struct {
	Documents []Document `json:"document"`
}

// Document struct
type Document struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

var app *App

func main() {
	app = NewApp()
	searchRes, searchErr := app.Search(app.Query, app.Source, app.Size, app.Offset)
	if searchErr != nil {
		panic(searchErr)
	}
	for _, d := range searchRes.Result.Documents {
		log.Println(d.Title)
	}
}

// NewApp constructor
func NewApp() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &App{
		URLTemplate: os.Getenv("URL_TEMPLATE"),
		Accept:      os.Getenv("ACCEPT_HEADER"),
		Query:       url.QueryEscape(os.Getenv("Q_PARAM")),
		Source:      url.QueryEscape(os.Getenv("SRC_PARAM")),
		Size:        url.QueryEscape(os.Getenv("SIZE_PARAM")),
		Offset:      url.QueryEscape(os.Getenv("OFFSET_PARAM")),
		Client:      http.Client{},
	}
}

// Search method
func (app *App) Search(q, src, size, offset string) (*Body, error) {
	url := fmt.Sprintf(
		app.URLTemplate,
		url.QueryEscape(q),
		url.QueryEscape(src),
		url.QueryEscape(size),
		url.QueryEscape(offset))

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Accept", app.Accept)

	response, responseErr := app.Client.Do(req)
	if responseErr != nil {
		return nil, responseErr
	}

	var body *Body
	decodeErr := json.NewDecoder(response.Body).Decode(&body)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return body, nil
}
