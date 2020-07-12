package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckin_should_be_ok(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(&Check{ID: 1, PlaceID: 212})
	req := httptest.NewRequest("POST", "http://www.example.com/req", payload)
	w := httptest.NewRecorder()

	var fn InFunc = func(ID, placeID int64) error {
		return nil
	}

	CheckIn(fn)(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(body)
}

func TestSealMiddleWare(t *testing.T) {
	payload := bytes.NewBuffer([]byte("eyJpZCIgOiAxMjMsICJwbGFjZUlEIjogMjIyfQ=="))
	req := httptest.NewRequest("POST", "http://www.example.com/req", payload)
	w := httptest.NewRecorder()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(b))
		w.Write(b)
	}

	SealMiddleWare()(http.HandlerFunc(handler)).ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(body)
}
