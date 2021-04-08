package pancors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandleProxy))
	defer ts.Close()

	type expected struct {
		statusCode int
		headers    map[string]string
	}

	tests := []struct {
		name     string
		url      string
		expected expected
	}{
		{
			"https url with params",
			"https://suggest.seznam.cz/slovnik/mix_cz_en?phrase=test&format=json-2",
			expected{
				http.StatusOK,
				map[string]string{
					"Access-Control-Allow-Origin":      "*",
					"Access-Control-Allow-Credentials": "true",
				},
			},
		},
		{
			"http url with params",
			"http://suggest.seznam.cz/slovnik/mix_cz_en?phrase=test&format=json-2",
			expected{
				http.StatusOK,
				map[string]string{
					"Access-Control-Allow-Origin":      "*",
					"Access-Control-Allow-Credentials": "true",
				},
			},
		},
		{
			"empty url",
			"",
			expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			"non http(s) url",
			"ftp://example.com",
			expected{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Logf("could not parse test server's url; got %v", err)
				t.FailNow()
			}

			q := u.Query()
			q.Set("url", tc.url)
			u.RawQuery = q.Encode()

			req, err := http.NewRequestWithContext(context.Background(), "GET", u.String(), nil)
			if err != nil {
				t.Logf("could not prepare a request; got %v", err)
				t.FailNow()
			}

			rsp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Log("could not fetch testing data")
				t.FailNow()
			}
			defer rsp.Body.Close()

			if rsp.StatusCode != tc.expected.statusCode {
				t.Logf("expected HTTP status code %d; got %d", tc.expected.statusCode, rsp.StatusCode)
				t.Fail()
			}

			for header, expected := range tc.expected.headers {
				actual := rsp.Header.Get(header)
				if actual != expected {
					t.Logf("expected header %s = %s; got: %v", header, expected, actual)
					t.Fail()
				}
			}
		})
	}
}
