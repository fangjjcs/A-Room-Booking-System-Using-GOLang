package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct{
	key string
	value string
}

var theTest = []struct{
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors", "GET", []postData{}, http.StatusOK},
	{"search-availability-get", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"eservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"search-availability-post", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-01-01"},{key: "end", value: "2021-01-02"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Mayor"},
		{key: "email", value: "john@mail.com"},
		{key: "phone", value: "0987-09889"},
	}, http.StatusOK},
}


func TestHandlers(t *testing.T){
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()
	for _, e := range theTest{
		if e.method == "GET" {
			res, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if res.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s, expected %d, but %d",e.name, e.expectedStatusCode, res.StatusCode)
			}
			
		} else{
			values := url.Values{}
			for _, x := range e.params{
				values.Add(x.key, x.value)
			}
			res, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if res.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s, expected %d, but %d",e.name, e.expectedStatusCode, res.StatusCode)
			}
		}
	}
}