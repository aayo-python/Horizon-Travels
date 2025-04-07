package controllers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequestMethod(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		allowedMethods []string
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "valid method",
			request:        &http.Request{Method: "GET"},
			allowedMethods: []string{"GET", "POST"},
			expectError:    false,
		},
		{
			name:           "invalid method",
			request:        &http.Request{Method: "DELETE"},
			allowedMethods: []string{"GET", "POST"},
			expectError:    true,
			errorMessage:   "method not allowed: DELETE",
		},
		{
			name:           "case insensitive method check",
			request:        &http.Request{Method: "get"},
			allowedMethods: []string{"GET", "POST"},
			expectError:    false,
		},
		{
			name:           "empty allowed methods",
			request:        &http.Request{Method: "GET"},
			allowedMethods: []string{},
			expectError:    true,
			errorMessage:   "method not allowed: GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequestMethod(tt.request, tt.allowedMethods...)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractParams(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedParams map[string]string
	}{
		{
			name: "valid URL parameters",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param2=value2",
				},
			},
			expectedParams: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
		},
		{
			name: "empty URL parameters",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "",
				},
			},
			expectedParams: map[string]string{},
		},
		{
			name: "URL parameters with empty values",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=&param2=value2",
				},
			},
			expectedParams: map[string]string{
				"param1": "",
				"param2": "value2",
			},
		},
		{
			name: "URL parameters with duplicate keys",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param1=value2",
				},
			},
			expectedParams: map[string]string{
				"param1": "value1", // First value is used
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := extractParams(tt.request)
			assert.Equal(t, tt.expectedParams, params)
		})
	}
}

func TestExtractParamValue(t *testing.T) {
	tests := []struct {
		name          string
		request       *http.Request
		paramName     string
		defaultValue  string
		expectedValue string
	}{
		{
			name: "parameter exists",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param2=value2",
				},
			},
			paramName:     "param1",
			defaultValue:  "default",
			expectedValue: "value1",
		},
		{
			name: "parameter does not exist",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param2=value2",
				},
			},
			paramName:     "param3",
			defaultValue:  "default",
			expectedValue: "default",
		},
		{
			name: "parameter exists but empty",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=&param2=value2",
				},
			},
			paramName:     "param1",
			defaultValue:  "default",
			expectedValue: "",
		},
		{
			name:          "nil request",
			request:       nil,
			paramName:     "param1",
			defaultValue:  "default",
			expectedValue: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := extractParamValue(tt.request, tt.paramName, tt.defaultValue)
			assert.Equal(t, tt.expectedValue, value)
		})
	}
}

func TestExtractParamValues(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		paramName      string
		expectedValues []string
	}{
		{
			name: "single parameter value",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param2=value2",
				},
			},
			paramName:      "param1",
			expectedValues: []string{"value1"},
		},
		{
			name: "multiple parameter values",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param1=value2&param1=value3",
				},
			},
			paramName:      "param1",
			expectedValues: []string{"value1", "value2", "value3"},
		},
		{
			name: "parameter does not exist",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=value1&param2=value2",
				},
			},
			paramName:      "param3",
			expectedValues: []string{},
		},
		{
			name: "parameter exists but empty",
			request: &http.Request{
				URL: &url.URL{
					RawQuery: "param1=&param2=value2",
				},
			},
			paramName:      "param1",
			expectedValues: []string{""},
		},
		{
			name:           "nil request",
			request:        nil,
			paramName:      "param1",
			expectedValues: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := extractParamValues(tt.request, tt.paramName)
			assert.Equal(t, tt.expectedValues, values)
		})
	}
}
