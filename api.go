package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/go-resty/resty/v2"
	//"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

const apiEndpoint = "https://api.openai.com/v1/chat/completions"

func (s *APIServer) Start(ctx context.Context) error {
	err := s.rdb.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to connect to redis")
	}

	fmt.Println("Redis server starting")

	ch := make(chan error, 1)

	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/newsfeed", makeHTTPHandleFunc(s.handleNewsFeed))

	defer func() {
		if err := s.rdb.Close(); err != nil {
			fmt.Println("Failed to close redis")
		} else {
			fmt.Println("Redis shutdown successful")
		}
	}()

	fmt.Println("starting server")

	server := &http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	go func() {
		log.Println("JSON API server running on port: ", s.listenAddr)
		err = server.ListenAndServe()

		if err != nil {
			ch <- fmt.Errorf("failed to start server %w", err)
		}

		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		fmt.Println("Done code hit.")
		timeout, _ := context.WithTimeout(context.Background(), time.Second*10)
		//defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(v)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func (s *APIServer) Init() {
	s.newsParser = newHTMLParser("https://www.psychologytoday.com/us/essentials")
	s.newsParser.InitCollector()
	s.newsParser.Collect()
}

func (s *APIServer) Run() {

}

func newHTMLParser(url string) *HTMLParser {
	return &HTMLParser{
		url: url,
	}
}

func (s *APIServer) handleNewsFeed(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetNews(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetNews(w http.ResponseWriter, r *http.Request) error {

	values := make([]*PageNode, 0, len(s.newsParser.pageNodes))

	for _, v := range s.newsParser.pageNodes {
		//fmt.Println("Key: " + key + " Value: " + v.title)
		values = append(values, v)
	}

	//newNode := NewPageNode()
	//account := NewAccount("Shawn", "McLean")

	return writeJSON(w, http.StatusOK, values)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Shawn", "McLean")

	return writeJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
