package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Fizzbuzz(t *testing.T) {
	req, err := http.NewRequest("GET", "/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	Fizzbuzz(res, req)
	exp := "fizz"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s got %s", exp, act)
	}
}
