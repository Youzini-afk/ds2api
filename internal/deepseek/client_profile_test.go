package deepseek

import (
	"testing"

	"ds2api/internal/config"
)

func TestAuthHeadersUsesCompatUpstreamWebProfile(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{"compat":{"preset":"shallowseek_compat"}}`)
	store := config.LoadStore()
	c := &Client{Store: store}

	headers := c.authHeaders("token-x")
	if headers["x-client-platform"] != "web" {
		t.Fatalf("expected web profile headers, got x-client-platform=%q", headers["x-client-platform"])
	}
	if headers["x-app-version"] == "" || headers["x-client-timezone-offset"] == "" {
		t.Fatalf("expected shallowseek-style web headers, got %#v", headers)
	}
	if headers["authorization"] != "Bearer token-x" {
		t.Fatalf("expected bearer authorization header, got %q", headers["authorization"])
	}
}

func TestAuthHeadersDefaultToAndroidWhenStoreMissing(t *testing.T) {
	c := &Client{}
	headers := c.authHeaders("token-y")
	if headers["x-client-platform"] != "android" {
		t.Fatalf("expected android default profile headers, got x-client-platform=%q", headers["x-client-platform"])
	}
}
