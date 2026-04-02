package openai

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ds2api/internal/util"
)

func TestHandleNonStreamHidesReasoningWhenExposeDisabled(t *testing.T) {
	h := &Handler{}
	resp := makeSSEHTTPResponse(
		`data: {"p":"response/thinking_content","v":"思考内容"}`,
		`data: {"p":"response/content","v":"正文输出"}`,
		`data: [DONE]`,
	)
	rec := httptest.NewRecorder()

	h.handleNonStream(rec, context.Background(), resp, "cid-hide", "deepseek-reasoner", "prompt", true, false, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
	out := decodeJSONBody(t, rec.Body.String())
	choices, _ := out["choices"].([]any)
	choice, _ := choices[0].(map[string]any)
	msg, _ := choice["message"].(map[string]any)
	if _, ok := msg["reasoning_content"]; ok {
		t.Fatalf("expected reasoning_content to be hidden, got %#v", msg["reasoning_content"])
	}
	if got, _ := msg["content"].(string); got == "" {
		t.Fatalf("expected normal content to remain, got %#v", msg["content"])
	}
}

func TestHandleStreamHidesReasoningDeltaWhenExposeDisabled(t *testing.T) {
	h := &Handler{}
	resp := makeSSEHTTPResponse(
		`data: {"p":"response/thinking_content","v":"思考中"}`,
		`data: {"p":"response/content","v":"最终答案"}`,
		`data: [DONE]`,
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	h.handleStream(rec, req, resp, "cid-stream-hide", "deepseek-reasoner", "prompt", true, false, false, nil)

	frames, _ := parseSSEDataFrames(t, rec.Body.String())
	for _, frame := range frames {
		choices, _ := frame["choices"].([]any)
		for _, item := range choices {
			choice, _ := item.(map[string]any)
			delta, _ := choice["delta"].(map[string]any)
			if _, ok := delta["reasoning_content"]; ok {
				t.Fatalf("expected no reasoning_content delta when expose_reasoning=false, body=%s", rec.Body.String())
			}
		}
	}
}

func TestHandleResponsesStreamHidesReasoningDeltaWhenExposeDisabled(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	rec := httptest.NewRecorder()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(
			`data: {"p":"response/thinking_content","v":"thought"}` + "\n" +
				`data: {"p":"response/content","v":"answer"}` + "\n" +
				`data: [DONE]` + "\n",
		)),
	}

	h.handleResponsesStream(rec, req, resp, "owner-a", "resp_hide", "deepseek-reasoner", "prompt", true, false, false, nil, util.DefaultToolChoicePolicy(), "")

	body := rec.Body.String()
	if strings.Contains(body, "event: response.reasoning.delta") {
		t.Fatalf("expected no response.reasoning.delta when expose_reasoning=false, body=%s", body)
	}
	if !strings.Contains(body, "event: response.completed") {
		t.Fatalf("expected response.completed event, body=%s", body)
	}
}
