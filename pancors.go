package pancors

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type CorsTransport http.Header

func (t CorsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", "*")
	res.Header.Set("Access-Control-Allow-Credentials", "true")

	return res, nil
}

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")

	urlParsed, err := url.Parse(urlParam)
	if err != nil || (urlParsed.Scheme != "http" && urlParsed.Scheme != "https") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = urlParsed
			r.Host = urlParsed.Host
		},
		Transport: CorsTransport(http.Header{}),
	}

	proxy.ServeHTTP(w, r)
}
