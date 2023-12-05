package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type OkTestHandler struct{}

func (h OkTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestContentTypeText(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        int
	}{
		{
			name:        "positive test",
			contentType: "text/plain",
			want:        200,
		},
		{
			name:        "negative test #1",
			contentType: "application/json",
			want:        400,
		},
	}

	handler := ContentTypeText(OkTestHandler{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requset := httptest.NewRequest("GET", "/test", nil)
			requset.Header.Add("Content-Type", test.contentType)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, requset)
			assert.Equal(t, test.want, w.Result().StatusCode)
		})
	}
}
