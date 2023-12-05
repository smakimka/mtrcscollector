package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodText(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   int
	}{
		{
			name:   "positive test",
			method: "POST",
			want:   200,
		},
		{
			name:   "negative test #1",
			method: "GET",
			want:   405,
		},
	}

	handler := MethodPOST(OkTestHandler{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requset := httptest.NewRequest(test.method, "/test", nil)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, requset)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
