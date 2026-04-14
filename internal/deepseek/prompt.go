package deepseek

import (
	"strings"

	"ds2api/internal/config"
	"ds2api/internal/prompt"
)

func MessagesPrepare(messages []map[string]any) string {
	return prompt.MessagesPrepare(messages)
}

func MessagesPrepareWithCompat(messages []map[string]any, resolvedModel, reasonerPromptMode string) string {
	useReasonerBoundary :=
		config.IsReasonerModel(strings.TrimSpace(resolvedModel)) &&
			strings.EqualFold(strings.TrimSpace(reasonerPromptMode), config.CompatReasonerPromptEndThink)
	return prompt.MessagesPrepareWithOptions(messages, prompt.PrepareOptions{
		ReasonerAssistantBoundary: useReasonerBoundary,
	})
}

func MessagesPrepareWithThinking(messages []map[string]any, thinkingEnabled bool) string {
	return prompt.MessagesPrepareWithThinking(messages, thinkingEnabled)
}
