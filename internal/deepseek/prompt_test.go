package deepseek

import (
	"strings"
	"testing"

	"ds2api/internal/config"
)

func TestMessagesPrepareWithCompatAddsReasonerBoundary(t *testing.T) {
	messages := []map[string]any{
		{"role": "assistant", "content": "历史回答"},
		{"role": "user", "content": "继续"},
	}
	got := MessagesPrepareWithCompat(messages, "deepseek-reasoner", config.CompatReasonerPromptEndThink)
	if !strings.Contains(got, "<｜Assistant｜><｜end▁of▁thinking｜>历史回答<｜end▁of▁sentence｜>") {
		t.Fatalf("expected reasoner assistant boundary with end_of_thinking, got %q", got)
	}
}

func TestMessagesPrepareWithCompatKeepsDefaultForNonReasoner(t *testing.T) {
	messages := []map[string]any{
		{"role": "assistant", "content": "历史回答"},
	}
	got := MessagesPrepareWithCompat(messages, "deepseek-chat", config.CompatReasonerPromptEndThink)
	if strings.Contains(got, "<｜end▁of▁thinking｜>") {
		t.Fatalf("did not expect end_of_thinking boundary for non-reasoner model, got %q", got)
	}
}
