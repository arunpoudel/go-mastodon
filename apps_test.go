package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterApp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if r.URL.Path != "/api/v1/apps" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if r.FormValue("redirect_uris") != "urn:ietf:wg:oauth:2.0:oob" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, `{"client_id": "foo", "client_secret": "bar"}`)
		return
	}))
	defer ts.Close()

	app, err := RegisterApp(context.Background(), &AppConfig{
		Server: ts.URL,
		Scopes: "read write follow",
	})
	if err != nil {
		t.Fatalf("should not be fail: %v", err)
	}
	if app.ClientID != "foo" {
		t.Fatalf("want %q but %q", "foo", app.ClientID)
	}
	if app.ClientSecret != "bar" {
		t.Fatalf("want %q but %q", "bar", app.ClientSecret)
	}
}
