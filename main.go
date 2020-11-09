package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	ID       string `json:"id"`
	Title    string `json:"Title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func articles(response http.ResponseWriter, request *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("methods", request.Method)
	if request.Method == "GET" {
		var allArticles []Article
		collection := client.Database("inshorts").Collection("news")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var article Article
			cursor.Decode(&article)
			allArticles = append(allArticles, article)
		}
		json.NewEncoder(response).Encode(allArticles)
	} else if request.Method == "POST" {

		var article Article
		_ = json.NewDecoder(request.Body).Decode(&article)
		collection := client.Database("inshorts").Collection("news")

		result, _ := collection.InsertOne(ctx, article)
		json.NewEncoder(response).Encode(result)
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))

	}
}

func search(response http.ResponseWriter, request *http.Request) {
	fmt.Println("GET params were:", request.URL.Query())
	query := request.URL.Query().Get("q")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	var article Article
	collection := client.Database("inshorts").Collection("news")
	fmt.Println(query)
	err = collection.FindOne(ctx, Article{Title: query, Subtitle: query, Content: query}).Decode(&article)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(article)
}

func getArticle(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Welcome to the getArticle page!")

	var param1, param2 = getCode(request, 267)

	fmt.Println("feshgfuheguhe")
	fmt.Println(param1)
	fmt.Println(param2)
}

func getCode(r *http.Request, defaultCode int) (int, string) {
	p := strings.Split(r.URL.Path, "/")
	fmt.Println(p[1])
	if len(p) == 1 {
		return defaultCode, p[0]
	} else if len(p) > 1 {
		code, err := strconv.Atoi(p[0])
		if err == nil {
			return code, p[1]
		} else {
			return defaultCode, p[1]
		}
	} else {
		return defaultCode, ""
	}
}

func handleRequests() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", articles)
	http.HandleFunc("/articles/search", search)
	http.HandleFunc("/articles/{id}", getArticle)

	http.ListenAndServe(":12345", nil)
}

func main() {

	handleRequests()
}
