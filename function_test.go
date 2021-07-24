package ddnsfunction

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		inputMethod     string
		inputPath       string
		inputBody       string
		outputBody      string
		outputStatsCode int
	}{
		{
			inputMethod: "GET",
			inputPath:   "/", inputBody: ``,
			outputBody:      "Dynamic DNS Service (0.0.0)",
			outputStatsCode: 200,
		}, {
			inputMethod:     "GET",
			inputPath:       "/bob",
			inputBody:       ``,
			outputBody:      "Not Found",
			outputStatsCode: 404,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.inputMethod, test.inputPath, strings.NewReader(test.inputBody))
		req.Header.Add("Content-Type", "application/text")

		rr := httptest.NewRecorder()
		Handler(rr, req)

		if got := rr.Body.String(); got != test.outputBody {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.inputBody, got, test.outputBody)
		}
	}
}

func TestObtainVersion(t *testing.T) {
	result := ObtainVersion()
	EqualsString(t, "0.0.0", result)
}

func EqualsString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf(
			"Expected %s, but got %s",
			expected,
			actual,
		)
	}
}
