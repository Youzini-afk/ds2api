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
	useShallowseekCompat := strings.EqualFold(strings.TrimSpace(reasonerPromptMode), config.CompatReasonerPromptEndThink)
	useReasonerBoundary :=
		config.IsReasonerModel(strings.TrimSpace(resolvedModel)) &&
			useShallowseekCompat
	return prompt.MessagesPrepareWithOptions(messages, prompt.PrepareOptions{
		ReasonerAssistantBoundary: useReasonerBoundary,
		ShallowseekCompat:         useShallowseekCompat,
	})
}
