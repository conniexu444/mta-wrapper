package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conniexu444/mta-wrapper/handlers"
)

func TestArrivalsHandler_MissingRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/arrivals", nil)
	rr := httptest.NewRecorder()

	handlers.ArrivalsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}
