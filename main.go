package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Article struct {
	ID       string `json:"id"`
	Title    string `json:"Title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

var Articles []Article

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func articles(response http.ResponseWriter, request *http.Request) {
	fmt.Println("methods", request.Method)
	if request.Method == "GET" {
		json.NewEncoder(response).Encode(Articles)
	} else if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		var article Article
		json.Unmarshal([]byte(body), &article)
		fmt.Println("Printing response")
		fmt.Println(article.ID)
		fmt.Println(article.Title)
		fmt.Println(article.Subtitle)
		fmt.Println(article.Content)
	}
}

func search(response http.ResponseWriter, request *http.Request) {
	fmt.Println("GET params were:", request.URL.Query())
	query := request.URL.Query().Get("q")
	fmt.Fprintf(response, "Welcome to the searchPage!")
	fmt.Println(query)
}

func getArticle(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to the getArticle page!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", articles)
	http.HandleFunc("/articles/search", search)
	http.HandleFunc("/article/{id}", getArticle)
	http.ListenAndServe(":12345", nil)
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())

	Articles = []Article{
		Article{ID: "1", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "2", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "3", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "4", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "5", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "6", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
		Article{ID: "7", Title: "Hello 2", Subtitle: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
