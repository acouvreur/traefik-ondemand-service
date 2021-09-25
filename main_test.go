package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/client"
)

type ScalerMock struct {
	isUp bool
}

func (s ScalerMock) IsUp(client *client.Client, name string) bool {
	return s.isUp
}

func (ScalerMock) ScaleUp(client *client.Client, name string, replicas *uint64) {}

func (ScalerMock) ScaleDown(client *client.Client, name string) {}

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
			contentType: "application/json",
		},
		{
			desc:        "service is started",
			status:      "started",
			scaler:      ScalerMock{isUp: true},
			statusCode:  http.StatusCreated,
			contentType: "application/json",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			t.Logf("IsUp: %t", test.scaler.isUp)
			request := httptest.NewRequest("GET", "/?name=whoami&timeout=5m", nil)
			responseRecorder := httptest.NewRecorder()

			onDemandHandler := onDemand(nil, test.scaler)
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
