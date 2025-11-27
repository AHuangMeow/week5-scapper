package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Pid  string `bson:"Pid" json:"Pid"`
	Name string `bson:"name" json:"name"`
}

var Client *mongo.Client
var UserCollection *mongo.Collection

func InitMongoDB(connectionString, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	UserCollection = Client.Database(dbName).Collection("users")
	log.Println("Connected to MongoDB!")
}

func GetUserByID(id int, cookieValue string) {
	req, err := http.NewRequest("GET", "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term="+strconv.Itoa(id)+"&_=1764144604145", nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: cookieValue,
	}
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(user)
	}
}

func Crawl() {
	InitMongoDB("mongodb://localhost:27017", "userdb")
	cookieValue := GetCookie()
	var wg sync.WaitGroup

	for i := 2024000000; i <= 2024999999; i += 50000 {
		wg.Add(1)
		go func(from int) {
			for id := from; id < from+50000; id++ {
				GetUserByID(id, cookieValue)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
