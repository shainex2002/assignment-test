package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDumpRequest(t *testing.T) {

	// Create a request to pass to handler.
	req, err := http.NewRequest("GET", "/test200", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := dumpRequest(http.StatusOK, "test200")

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
