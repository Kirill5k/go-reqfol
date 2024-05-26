package proxy

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRequestMetadata_String(t *testing.T) {
	tests := []struct {
		name                         string
		request                      RequestMetadata
		expectedStringRepresentation string
	}{
		{
			name: "Convert request metadata to string",
			request: RequestMetadata{
				Method: "GET",
				Url:    "https://google.com",
				Headers: map[string]string{
					"Header1": "Value1",
					"Header2": "Value2",
				},
				QueryParams: map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
			expectedStringRepresentation: "GET https://google.com?param1=value1&param2=value2 {Header1:Value1, Header2:Value2}",
		},
		{
			name: "Append query params to string when original url already contains other parameters",
			request: RequestMetadata{
				Method: "GET",
				Url:    "https://google.com?foo=bar",
				QueryParams: map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
			expectedStringRepresentation: "GET https://google.com?foo=bar&param1=value1&param2=value2 {}",
		},
		{
			name: "Handle simple cases without query params",
			request: RequestMetadata{
				Method: "GET",
				Url:    "https://google.com?foo=bar",
			},
			expectedStringRepresentation: "GET https://google.com?foo=bar {}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.request.String()

			require.Exactly(t, test.expectedStringRepresentation, result, test.name)
		})
	}
}
