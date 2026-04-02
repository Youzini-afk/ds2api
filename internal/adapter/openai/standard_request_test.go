package openai

import (
	"strings"
	"testing"

	"ds2api/internal/config"
	"ds2api/internal/util"
)

func newEmptyStoreForNormalizeTest(t *testing.T) *config.Store {
	t.Helper()
	t.Setenv("DS2API_CONFIG_JSON", `{}`)
	return config.LoadStore()
}

func TestNormalizeOpenAIChatRequest(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "gpt-5-codex",
		"messages": []any{
			map[string]any{"role": "user", "content": "hello"},
		},
		"temperature": 0.3,
		"stream":      true,
	}
	n, err := normalizeOpenAIChatRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ResolvedModel != "deepseek-reasoner" {
		t.Fatalf("unexpected resolved model: %s", n.ResolvedModel)
	}
	if !n.Stream {
		t.Fatalf("expected stream=true")
	}
	if _, ok := n.PassThrough["temperature"]; !ok {
		t.Fatalf("expected temperature passthrough")
	}
	if n.FinalPrompt == "" {
		t.Fatalf("expected non-empty final prompt")
	}
}

func TestNormalizeOpenAIResponsesRequestInput(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model":        "gpt-4o",
		"input":        "ping",
		"instructions": "system",
	}
	n, err := normalizeOpenAIResponsesRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ResolvedModel != "deepseek-chat" {
		t.Fatalf("unexpected resolved model: %s", n.ResolvedModel)
	}
	if len(n.Messages) != 2 {
		t.Fatalf("expected 2 normalized messages, got %d", len(n.Messages))
	}
}

func TestNormalizeOpenAIResponsesRequestToolChoiceRequired(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "gpt-4o",
		"input": "ping",
		"tools": []any{
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": "search",
					"parameters": map[string]any{
						"type": "object",
					},
				},
			},
		},
		"tool_choice": "required",
	}
	n, err := normalizeOpenAIResponsesRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ToolChoice.Mode != util.ToolChoiceRequired {
		t.Fatalf("expected tool choice mode required, got %q", n.ToolChoice.Mode)
	}
	if len(n.ToolNames) != 1 || n.ToolNames[0] != "search" {
		t.Fatalf("unexpected tool names: %#v", n.ToolNames)
	}
}

func TestNormalizeOpenAIResponsesRequestToolChoiceForcedFunction(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "gpt-4o",
		"input": "ping",
		"tools": []any{
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": "search",
				},
			},
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": "read_file",
				},
			},
		},
		"tool_choice": map[string]any{
			"type": "function",
			"name": "read_file",
		},
	}
	n, err := normalizeOpenAIResponsesRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ToolChoice.Mode != util.ToolChoiceForced {
		t.Fatalf("expected tool choice mode forced, got %q", n.ToolChoice.Mode)
	}
	if n.ToolChoice.ForcedName != "read_file" {
		t.Fatalf("expected forced tool name read_file, got %q", n.ToolChoice.ForcedName)
	}
	if len(n.ToolNames) != 1 || n.ToolNames[0] != "read_file" {
		t.Fatalf("expected filtered tool names [read_file], got %#v", n.ToolNames)
	}
}

func TestNormalizeOpenAIResponsesRequestToolChoiceForcedUndeclaredFails(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "gpt-4o",
		"input": "ping",
		"tools": []any{
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": "search",
				},
			},
		},
		"tool_choice": map[string]any{
			"type": "function",
			"name": "read_file",
		},
	}
	if _, err := normalizeOpenAIResponsesRequest(store, req, ""); err == nil {
		t.Fatalf("expected forced undeclared tool to fail")
	}
}

func TestNormalizeOpenAIResponsesRequestToolChoiceNoneKeepsToolDetectionEnabled(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "gpt-4o",
		"input": "ping",
		"tools": []any{
			map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": "search",
				},
			},
		},
		"tool_choice": "none",
	}
	n, err := normalizeOpenAIResponsesRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ToolChoice.Mode != util.ToolChoiceNone {
		t.Fatalf("expected tool choice mode none, got %q", n.ToolChoice.Mode)
	}
	if len(n.ToolNames) == 0 {
		t.Fatalf("expected tool detection sentinel when tool_choice=none, got %#v", n.ToolNames)
	}
}

func TestNormalizeOpenAIChatRequestExposeReasoningAlwaysByDefault(t *testing.T) {
	store := newEmptyStoreForNormalizeTest(t)
	req := map[string]any{
		"model": "deepseek-reasoner",
		"messages": []any{
			map[string]any{"role": "user", "content": "hello"},
		},
	}
	n, err := normalizeOpenAIChatRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if !n.ExposeReasoning {
		t.Fatalf("expected expose_reasoning=true in always mode")
	}
}

func TestNormalizeOpenAIChatRequestExposeReasoningRequestOptIn(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{"compat":{"preset":"shallowseek_compat"}}`)
	store := config.LoadStore()
	req := map[string]any{
		"model": "deepseek-reasoner",
		"messages": []any{
			map[string]any{"role": "user", "content": "hello"},
		},
	}

	n, err := normalizeOpenAIChatRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if n.ExposeReasoning {
		t.Fatalf("expected expose_reasoning=false when include_reasoning is omitted in request_opt_in mode")
	}

	req["include_reasoning"] = true
	n2, err := normalizeOpenAIChatRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if !n2.ExposeReasoning {
		t.Fatalf("expected expose_reasoning=true when include_reasoning=true in request_opt_in mode")
	}
}

func TestNormalizeOpenAIChatRequestShallowseekPresetKeepsDefaultPromptStyle(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{"compat":{"preset":"shallowseek_compat"}}`)
	store := config.LoadStore()
	req := map[string]any{
		"model": "deepseek-reasoner",
		"messages": []any{
			map[string]any{"role": "system", "content": "你是助手"},
			map[string]any{"role": "user", "content": "你好"},
			map[string]any{"role": "assistant", "content": "你好呀"},
			map[string]any{"role": "user", "content": "继续"},
		},
	}
	n, err := normalizeOpenAIChatRequest(store, req, "")
	if err != nil {
		t.Fatalf("normalize failed: %v", err)
	}
	if !strings.HasPrefix(n.FinalPrompt, "<system_instructions>\n你是助手\n</system_instructions>\n\n<｜User｜>你好") {
		t.Fatalf("expected default prompt prefix, got %q", n.FinalPrompt)
	}
	if !strings.Contains(n.FinalPrompt, "<｜Assistant｜><｜end▁of▁thinking｜>你好呀<｜end▁of▁sentence｜><｜User｜>继续") {
		t.Fatalf("expected only reasoner assistant boundary override, got %q", n.FinalPrompt)
	}
}
