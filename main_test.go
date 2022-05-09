package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acouvreur/tinykv"
)

type ScalerMock struct {
	isUp bool
}

func (s ScalerMock) IsUp(name string) bool {
	return s.isUp
}

func (ScalerMock) ScaleUp(name string) error { return nil }

func (ScalerMock) ScaleDown(name string) error { return nil }

func TestOndemand_ServeHTTP(t *testing.T) {
	testCases := []struct {
		desc        string
		scaler      ScalerMock
		status      string
		statusCode  int
		contentType string
	}{
		{
			desc:        "service is starting",
			status:      "starting",
			scaler:      ScalerMock{isUp: false},
			statusCode:  http.StatusAccepted,
			contentType: "text/plain",
		},
		{
			desc:        "service is started",
			status:      "started",
			scaler:      ScalerMock{isUp: true},
			statusCode:  http.StatusCreated,
			contentType: "text/plain",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			t.Logf("IsUp: %t", test.scaler.isUp)
			request := httptest.NewRequest("GET", "/?name=whoami&timeout=5m", nil)
			responseRecorder := httptest.NewRecorder()

			store := tinykv.New[OnDemandRequestState](time.Second * 20)

			onDemandHandler := onDemand(test.scaler, store)
			onDemandHandler(responseRecorder, request)

			body := responseRecorder.Body.String()

			if responseRecorder.Code != test.statusCode {
				t.Errorf("Want status '%d', got '%d'", test.statusCode, responseRecorder.Code)
			}

			if responseRecorder.Body.String() != test.status {
				t.Errorf("Want body '%s', got '%s'", test.status, body)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Errorf("Want content type '%s', got '%s'", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}
		})
	}
}
