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

func TestMessagesPrepareWithCompatKeepsDefaultPromptStyleWhenBoundaryEnabled(t *testing.T) {
	messages := []map[string]any{
		{"role": "system", "content": "你是助手"},
		{"role": "user", "content": "你好"},
		{"role": "assistant", "content": "你好呀"},
		{"role": "user", "content": "继续"},
	}
	got := MessagesPrepareWithCompat(messages, "deepseek-reasoner", config.CompatReasonerPromptEndThink)
	if !strings.HasPrefix(got, "<system_instructions>\n你是助手\n</system_instructions>\n\n<｜User｜>你好") {
		t.Fatalf("expected default system/user formatting, got %q", got)
	}
	if !strings.Contains(got, "<｜Assistant｜><｜end▁of▁thinking｜>你好呀<｜end▁of▁sentence｜><｜User｜>继续") {
		t.Fatalf("expected only reasoner assistant boundary to change, got %q", got)
	}
}
