package deepseek

import "testing"

func TestSharedConstantsLoaded(t *testing.T) {
	if BaseHeaders["x-client-platform"] != "android" {
		t.Fatalf("unexpected base header x-client-platform=%q", BaseHeaders["x-client-platform"])
	}
	if len(SkipContainsPatterns) == 0 {
		t.Fatal("expected skip contains patterns to be loaded")
	}
	if _, ok := SkipExactPathSet["response/search_status"]; !ok {
		t.Fatal("expected response/search_status in exact skip path set")
	}
}

func TestBaseHeadersForProfileWeb(t *testing.T) {
	web := BaseHeadersForProfile("web")
	if web["x-client-platform"] != "web" {
		t.Fatalf("expected web profile x-client-platform=web, got %q", web["x-client-platform"])
	}
	if web["Origin"] == "" || web["Referer"] == "" {
		t.Fatalf("expected web profile Origin/Referer, got %#v", web)
	}
}

func TestBaseHeadersForProfileDefaultAndroid(t *testing.T) {
	android := BaseHeadersForProfile("unknown")
	if android["x-client-platform"] != "android" {
		t.Fatalf("expected android fallback profile, got %q", android["x-client-platform"])
	}
}
