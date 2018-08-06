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

type countData struct {
	Data string `json:"Data"`
}

type execCommand struct {
	Command string `json:"Command"`
}

// Fizzbuzz returns fizz if num dividable with 3 buzz if dividable with 5 and fizzbuzz if dividable with both.
func Fizzbuzz(res http.ResponseWriter, req *http.Request) {
	msg := ""
	params := mux.Vars(req)
	myInt, err := strconv.Atoi(params["num"])
	if err != nil {
		json.NewEncoder(res).Encode(params["num"] + " is not a number")
		// json.NewEncoder(res).Encode(params)
		return
	}
	if myInt%3 == 0 {
		msg += "fizz"
	}
	if myInt%5 == 0 {
		msg += "buzz"
	}
	json.NewEncoder(res).Encode(msg)
}

// Count counts the occurance of individual characters in Req.Body.Data.
func Count(res http.ResponseWriter, req *http.Request) {
	var myData countData
	var myMap = make(map[string]int)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&myData)
	if err != nil {
		json.NewEncoder(res).Encode("Invalid data")
		return
	}
	for i := range myData.Data {
		char := string(myData.Data[i])
		if myMap[char] == 0 {
			myMap[char] = 1
		} else {
			myMap[char]++
		}
	}
	json.NewEncoder(res).Encode(myMap)
}

// Exec executes a shell command on the host and sends back the output.
func Exec(res http.ResponseWriter, req *http.Request) {
	var command execCommand
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&command)
	if err != nil {
		json.NewEncoder(res).Encode("Invalid Command")
		return
	}
	out, err := exec.Command("sh", "-c", command.Command).Output()
	if err != nil {
		json.NewEncoder(res).Encode("Error during command execution")
		return
	}
	json.NewEncoder(res).Encode(string(out))
}

func custom404() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json.NewEncoder(res).Encode("Wrong way 404 🐸")
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
