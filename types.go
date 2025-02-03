package main

import (
	//"context"
	//"fmt"
	"math/rand"

	"github.com/gocolly/colly"
	"github.com/redis/go-redis/v9"
)

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

type APIServer struct {
	listenAddr string
	newsParser *HTMLParser
	rdb        *redis.Client
}

func newAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		rdb:        redis.NewClient(&redis.Options{}),
	}
}

type PageNode struct {
	Title string `json:"title"`
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

func NewPageNode() *PageNode {
	return &PageNode{
		Title: "Shawn",
		Topic: "McLean",
		Body:  "sux",
	}
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
		Balance:   0,
	}
}

type HTMLParser struct {
	url       string
	c         *colly.Collector
	pageNodes map[string]*PageNode
	hit       int
}
