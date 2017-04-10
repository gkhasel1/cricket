package main

import (
	"testing"
	"net/http"
    "net/http/httptest"
)

/**
 * This is by no means an exhaustive test suite but does show some golang
 * testing mechanisms. Ideally we would mock our the ES database and test
 * our db operations and handlers more thoroughly.
 **/

func TestIndexHandler(t *testing.T) {
	/**
	 * Test our index handler responds appropriately
	 **/
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder() // This is like a responsewriter but for testing.
    handler := http.HandlerFunc(IndexHandler)
    handler.ServeHTTP(rr, req) // call the handler HTTP serve directly

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "CRICKET\na metrics server written in go, backed by elasticseach\n"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestHealthHandler(t *testing.T) {
	/**
	 * Test that our health check returns status: 200 OK
	 **/
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder() // This is like a responsewriter but for testing.
    handler := http.HandlerFunc(HealthHandler)
    handler.ServeHTTP(rr, req) // call the handler HTTP serve directly

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "OK"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestGetStatsHandlerNoParams(t *testing.T) {
	/**
	 * Test that GetStats fails with no query parameters
	 **/
    req, err := http.NewRequest("GET", "/GetStats", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder() // This is like a responsewriter but for testing.
    handler := http.HandlerFunc(GetStatsHandler)
    handler.ServeHTTP(rr, req) // call the handler HTTP serve directly

    if status := rr.Code; status != 500 {
        t.Errorf("handler returned wrong status code: got %v want %v", status, 500)
    }
}
