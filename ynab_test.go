package ynab

import (
	"context"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	//"strings"
	"testing"
	//"time"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server

	accessToken = os.Getenv("YNAB_ACCESS_TOKEN")
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewDefaultClient(accessToken)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parseForm(): %v", err)
	}

	if !reflect.DeepEqual(expected, r.Form) {
		t.Errorf("Request parameters = %v, expected %v", r.Form, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		//"Account",
		//"Actions",
		//"Domains",
		//"Droplets",
		//"DropletActions",
		//"Images",
		//"ImageActions",
		//"Keys",
		//"Regions",
		//"Sizes",
		//"FloatingIPs",
		//"FloatingIPActions",
		//"Tags",
		"UserService",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		if cv.FieldByName(s).IsNil() {
			t.Errorf("c.%s shouldn't be nil", s)
		}
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	if c.BaseURL == nil || c.BaseURL.String() != DefaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, DefaultBaseURL)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultBaseURL(t, c)
	testClientServices(t, c)
}

func TestNewDefaultClient(t *testing.T) {
	c := NewDefaultClient(accessToken)
	testClientDefaults(t, c)
}

func TestNewRequest(t *testing.T) {
	c := NewDefaultClient(accessToken)

	inURL, outURL := "foo", DefaultBaseURL+"foo"
	inBody, outBody := &UserResponse{Data: UserWrapper{User: User{Id: "2590e617-2d1c-42bb-a105-42888579366d"}}},
		`{"data":{"user":{"id":"2590e617-2d1c-42bb-a105-42888579366d"}}}`+"\n"
	req, _ := c.newRequest(http.MethodGet, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("newRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("newRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewDefaultClient(accessToken)
	_, err := c.newRequest(http.MethodGet, ":", nil)
	testURLParseError(t, err)
}
