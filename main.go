package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
)

type resMsg struct {
	Body string `json:"resMsg,omitempty"`
}

type countData struct {
	Data string `json:"Data"`
}

type execCommand struct {
	Command string `json:"Command"`
}

// Fizzbuzz returns fizz if num dividable with 3 buzz if dividable with 5 and fizzbuzz if dividable with both.
func Fizzbuzz(wr http.ResponseWriter, req *http.Request) {
	msg := ""
	params := mux.Vars(req)
	myInt, err := strconv.Atoi(params["num"])
	if err != nil {
		res := resMsg{params["num"] + " is not a number"}
		json.NewEncoder(wr).Encode(res)
		return
	}
	if myInt%3 == 0 {
		msg += "fizz"
	}
	if myInt%5 == 0 {
		msg += "buzz"
	}
	res := resMsg{msg}
	json.NewEncoder(wr).Encode(res.Body)
}

// Count counts the occurance of individual characters in Req.Body.Data.
func Count(wr http.ResponseWriter, req *http.Request) {
	var myData countData
	var myMap = make(map[string]int)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&myData)
	if err != nil {
		res := resMsg{"Invalid data"}
		json.NewEncoder(wr).Encode(res)
		return
	}
	for i, _ := range myData.Data {
		char := string(myData.Data[i])
		if myMap[char] == 0 {
			myMap[char] = 1
		} else {
			myMap[char]++
		}
	}
	json.NewEncoder(wr).Encode(myMap)
}

// Exec executes a shell command on the host and sends back the output.
func Exec(wr http.ResponseWriter, req *http.Request) {
	var command execCommand
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&command)
	if err != nil {
		res := resMsg{"Invalid Command"}
		json.NewEncoder(wr).Encode(res)
		return
	}
	out, err := exec.Command("sh", "-c", command.Command).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("OUT : %s", out)
	res := resMsg{string(out)}
	json.NewEncoder(wr).Encode(res.Body)
}

func custom404() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		json.NewEncoder(wr).Encode("Wrong way 404 🐸")
	})
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = custom404()
	router.HandleFunc("/fizzbuzz/{num}", Fizzbuzz).Methods("GET")
	router.HandleFunc("/count", Count).Methods("POST")
	router.HandleFunc("/exec", Exec).Methods("POST")
	fmt.Println("Fizzbuzz server listening on PORT 8000")
	log.Fatal(http.ListenAndServe(":8000", logger.Handler(router, os.Stdout, logger.CommonLoggerType)))
}
