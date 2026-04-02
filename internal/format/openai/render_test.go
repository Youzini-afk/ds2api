package openai

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildResponseObjectToolCallsFollowChatShape(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"",
		`{"tool_calls":[{"name":"search","input":{"q":"golang"}}]}`,
		[]string{"search"},
		true,
	)

	outputText, _ := obj["output_text"].(string)
	if outputText != "" {
		t.Fatalf("expected output_text to be hidden for tool calls, got %q", outputText)
	}

	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected function_call output only, got %#v", obj["output"])
	}

	first, _ := output[0].(map[string]any)
	if first["type"] != "function_call" {
		t.Fatalf("expected first output item type function_call, got %#v", first["type"])
	}
	if first["call_id"] == "" {
		t.Fatalf("expected function_call item to have call_id, got %#v", first)
	}
	if first["name"] != "search" {
		t.Fatalf("unexpected function name: %#v", first["name"])
	}
	argsRaw, _ := first["arguments"].(string)
	var args map[string]any
	if err := json.Unmarshal([]byte(argsRaw), &args); err != nil {
		t.Fatalf("arguments should be valid json string, got=%q err=%v", argsRaw, err)
	}
	if args["q"] != "golang" {
		t.Fatalf("unexpected arguments: %#v", args)
	}
}

func TestBuildResponseObjectPromotesMixedProseToolPayloadToFunctionCall(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"",
		`示例格式：{"tool_calls":[{"name":"search","input":{"q":"golang"}}]}，但这条是普通回答。`,
		[]string{"search"},
		true,
	)

	outputText, _ := obj["output_text"].(string)
	if outputText != "" {
		t.Fatalf("expected output_text hidden for mixed prose tool payload, got %q", outputText)
	}
	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected one function_call output item, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "function_call" {
		t.Fatalf("expected function_call output type, got %#v", first["type"])
	}
}

func TestBuildResponseObjectKeepsFencedToolPayloadAsText(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"",
		"```json\n{\"tool_calls\":[{\"name\":\"search\",\"input\":{\"q\":\"golang\"}}]}\n```",
		[]string{"search"},
		true,
	)

	outputText, _ := obj["output_text"].(string)
	if !strings.Contains(outputText, "\"tool_calls\"") {
		t.Fatalf("expected output_text to preserve fenced tool payload, got %q", outputText)
	}
	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected one message output item, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "message" {
		t.Fatalf("expected message output type, got %#v", first["type"])
	}
}

// Backward-compatible alias for historical test name used in CI logs.
func TestBuildResponseObjectPromotesFencedToolPayloadToFunctionCall(t *testing.T) {
	TestBuildResponseObjectKeepsFencedToolPayloadAsText(t)
}

func TestBuildResponseObjectReasoningOnlyFallsBackToOutputText(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"internal thinking content",
		"",
		nil,
		true,
	)

	outputText, _ := obj["output_text"].(string)
	if outputText == "" {
		t.Fatalf("expected output_text fallback from reasoning when final text is empty")
	}

	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected one output item, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "message" {
		t.Fatalf("expected output type message, got %#v", first["type"])
	}
	content, _ := first["content"].([]any)
	if len(content) == 0 {
		t.Fatalf("expected reasoning content, got %#v", first["content"])
	}
	block0, _ := content[0].(map[string]any)
	if block0["type"] != "reasoning" {
		t.Fatalf("expected first content block reasoning, got %#v", block0["type"])
	}
}

func TestBuildResponseObjectKeepsReasoningWhenToolCallsArePresent(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"internal thinking content",
		`{"tool_calls":[{"name":"search","input":{"q":"golang"}}]}`,
		[]string{"search"},
		true,
	)

	outputText, _ := obj["output_text"].(string)
	if outputText != "" {
		t.Fatalf("expected output_text to stay hidden for tool calls, got %q", outputText)
	}

	output, _ := obj["output"].([]any)
	if len(output) != 2 {
		t.Fatalf("expected reasoning message + function_call, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "message" {
		t.Fatalf("expected first output item message, got %#v", first["type"])
	}
	content, _ := first["content"].([]any)
	if len(content) != 1 {
		t.Fatalf("expected one reasoning content block, got %#v", first["content"])
	}
	block0, _ := content[0].(map[string]any)
	if block0["type"] != "reasoning" {
		t.Fatalf("expected reasoning content block, got %#v", block0["type"])
	}
	second, _ := output[1].(map[string]any)
	if second["type"] != "function_call" {
		t.Fatalf("expected second output item function_call, got %#v", second["type"])
	}
}

func TestBuildResponseObjectHidesReasoningWhenToolCallsArePresentButExposureDisabled(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		"internal thinking content",
		`{"tool_calls":[{"name":"search","input":{"q":"golang"}}]}`,
		[]string{"search"},
		false,
	)

	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected function_call output only when reasoning is hidden, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "function_call" {
		t.Fatalf("expected function_call output type, got %#v", first["type"])
	}
}

func TestBuildResponseObjectIgnoresToolCallFromThinkingChannel(t *testing.T) {
	obj := BuildResponseObject(
		"resp_test",
		"gpt-4o",
		"prompt",
		`{"tool_calls":[{"name":"search","input":{"q":"from-thinking"}}]}`,
		"",
		[]string{"search"},
		true,
	)

	output, _ := obj["output"].([]any)
	if len(output) != 1 {
		t.Fatalf("expected one message output item, got %#v", obj["output"])
	}
	first, _ := output[0].(map[string]any)
	if first["type"] != "message" {
		t.Fatalf("expected output message, got %#v", first["type"])
	}
}
