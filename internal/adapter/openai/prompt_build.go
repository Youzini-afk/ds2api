package openai

import (
	"ds2api/internal/deepseek"
	"ds2api/internal/util"
)

func buildOpenAIFinalPrompt(messagesRaw []any, toolsRaw any, traceID string, thinkingEnabled bool, resolvedModel, reasonerPromptMode string) (string, []string) {
	return buildOpenAIFinalPromptWithPolicy(messagesRaw, toolsRaw, traceID, util.DefaultToolChoicePolicy(), thinkingEnabled, resolvedModel, reasonerPromptMode)
}

func buildOpenAIFinalPromptWithPolicy(messagesRaw []any, toolsRaw any, traceID string, toolPolicy util.ToolChoicePolicy, thinkingEnabled bool, resolvedModel, reasonerPromptMode string) (string, []string) {
	messages := normalizeOpenAIMessagesForPrompt(messagesRaw, traceID)
	toolNames := []string{}
	if tools, ok := toolsRaw.([]any); ok && len(tools) > 0 {
		messages, toolNames = injectToolPrompt(messages, tools, toolPolicy)
	}
	_ = thinkingEnabled
	return deepseek.MessagesPrepareWithCompat(messages, resolvedModel, reasonerPromptMode), toolNames
}

// BuildPromptForAdapter exposes the OpenAI-compatible prompt building flow so
// other protocol adapters (for example Gemini) can reuse the same tool/history
// normalization logic and remain behavior-compatible with chat/completions.
func BuildPromptForAdapter(messagesRaw []any, toolsRaw any, traceID string, thinkingEnabled bool, resolvedModel, reasonerPromptMode string) (string, []string) {
	return buildOpenAIFinalPrompt(messagesRaw, toolsRaw, traceID, thinkingEnabled, resolvedModel, reasonerPromptMode)
}
