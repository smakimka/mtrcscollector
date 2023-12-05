package models

import (
	"net/http"
	"testing"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	type want struct {
		data MetricsUpdateData
		code int
	}
	tests := []struct {
		name    string
		path    string
		wantErr bool
		want    want
	}{
		{
			name:    "positive test #1",
			path:    "gauge/test/1.0",
			wantErr: false,
			want: want{
				data: MetricsUpdateData{
					Kind:  mtrcs.Gauge,
					Name:  "test",
					Value: 1.0,
				},
				code: http.StatusOK,
			},
		},
		{
			name:    "positive test #2",
			path:    "counter/test/1",
			wantErr: false,
			want: want{
				data: MetricsUpdateData{
					Kind:  mtrcs.Counter,
					Name:  "test",
					Value: 1,
				},
				code: http.StatusOK,
			},
		},
		{
			name:    "wrong metric value",
			path:    "counter/test/1.5",
			wantErr: true,
			want: want{
				data: MetricsUpdateData{
					Kind: mtrcs.Counter,
					Name: "test",
				},
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "wrong metric type value",
			path:    "not_existing/test/1",
			wantErr: true,
			want: want{
				data: MetricsUpdateData{},
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "not enough params #1",
			path:    "gauge/test",
			wantErr: true,
			want: want{
				data: MetricsUpdateData{},
				code: http.StatusNotFound,
			},
		},
		{
			name:    "not enough params #2",
			path:    "counter",
			wantErr: true,
			want: want{
				data: MetricsUpdateData{},
				code: http.StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := MetricsUpdateData{}
			code, err := data.ParsePath(test.path)

			assert.Equal(t, test.want.data, data)
			assert.Equal(t, test.want.code, code)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
