package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_handleError(t *testing.T) {
	type args struct {
		err        error
		statusCode int
		req        *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantBody errorResponse
	}{
		{
			name: "should return error response 400",
			args: args{
				err:        nil,
				statusCode: http.StatusBadRequest,
				req:        httptest.NewRequest(http.MethodGet, "/api/v1/test", nil),
			},
			wantBody: errorResponse{
				Error:      "Bad Request: 400",
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "should return error response 500",
			args: args{
				err:        nil,
				statusCode: http.StatusInternalServerError,
				req:        httptest.NewRequest(http.MethodGet, "/api/v1/test", nil),
			},
			wantBody: errorResponse{
				Error:      "Internal Server Error: 500",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handleError(tt.args.req.Context(), w, tt.args.err, tt.args.statusCode, false)

			if w.Code != tt.args.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.args.statusCode, w.Code)
			}

			// Decode the response body into struct for comparison
			var resBody errorResponse
			if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Errorf("handleError() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(resBody, tt.wantBody) {
				t.Errorf("handleError() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
