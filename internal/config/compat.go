package config

import "strings"

const (
	CompatPresetDefault          = "default"
	CompatPresetShallowseek      = "shallowseek_compat"
	CompatReasonerPromptDefault  = "default"
	CompatReasonerPromptEndThink = "end_of_thinking"
	CompatReasoningAlways        = "always"
	CompatReasoningRequestOptIn  = "request_opt_in"
	CompatUpstreamAndroid        = "android"
	CompatUpstreamWeb            = "web"
)

type EffectiveCompat struct {
	Preset             string
	ReasonerPromptMode string
	ReasoningExposure  string
	UpstreamProfile    string
}

func ResolveCompat(cfg CompatConfig) EffectiveCompat {
	preset := normalizeCompatPreset(cfg.Preset)
	out := EffectiveCompat{
		Preset:             preset,
		ReasonerPromptMode: CompatReasonerPromptDefault,
		ReasoningExposure:  CompatReasoningAlways,
		UpstreamProfile:    CompatUpstreamAndroid,
	}

	if preset == CompatPresetShallowseek {
		out.ReasonerPromptMode = CompatReasonerPromptEndThink
		out.ReasoningExposure = CompatReasoningRequestOptIn
		out.UpstreamProfile = CompatUpstreamWeb
	}

	if override := normalizeCompatReasonerPromptMode(cfg.ReasonerPromptModeOverride); override != "" {
		out.ReasonerPromptMode = override
	}
	if override := normalizeCompatReasoningExposure(cfg.ReasoningExposureOverride); override != "" {
		out.ReasoningExposure = override
	}
	if override := normalizeCompatUpstreamProfile(cfg.UpstreamProfileOverride); override != "" {
		out.UpstreamProfile = override
	}

	return out
}

func normalizeCompatPreset(raw string) string {
	switch normalizeCompatToken(raw) {
	case "", CompatPresetDefault:
		return CompatPresetDefault
	case CompatPresetShallowseek:
		return CompatPresetShallowseek
	default:
		return CompatPresetDefault
	}
}

func normalizeCompatReasonerPromptMode(raw string) string {
	switch normalizeCompatToken(raw) {
	case CompatReasonerPromptDefault:
		return CompatReasonerPromptDefault
	case CompatReasonerPromptEndThink:
		return CompatReasonerPromptEndThink
	default:
		return ""
	}
}

func normalizeCompatReasoningExposure(raw string) string {
	switch normalizeCompatToken(raw) {
	case CompatReasoningAlways:
		return CompatReasoningAlways
	case CompatReasoningRequestOptIn:
		return CompatReasoningRequestOptIn
	default:
		return ""
	}
}

func normalizeCompatUpstreamProfile(raw string) string {
	switch normalizeCompatToken(raw) {
	case CompatUpstreamAndroid:
		return CompatUpstreamAndroid
	case CompatUpstreamWeb:
		return CompatUpstreamWeb
	default:
		return ""
	}
}

func IsValidCompatPreset(raw string) bool {
	v := normalizeCompatToken(raw)
	return v == CompatPresetDefault || v == CompatPresetShallowseek
}

func IsValidCompatReasonerPromptMode(raw string) bool {
	v := normalizeCompatToken(raw)
	return v == CompatReasonerPromptDefault || v == CompatReasonerPromptEndThink
}

func IsValidCompatReasoningExposure(raw string) bool {
	v := normalizeCompatToken(raw)
	return v == CompatReasoningAlways || v == CompatReasoningRequestOptIn
}

func IsValidCompatUpstreamProfile(raw string) bool {
	v := normalizeCompatToken(raw)
	return v == CompatUpstreamAndroid || v == CompatUpstreamWeb
}

func compatHasValues(cfg CompatConfig) bool {
	return cfg.WideInputStrictOutput != nil ||
		strings.TrimSpace(cfg.Preset) != "" ||
		strings.TrimSpace(cfg.ReasonerPromptModeOverride) != "" ||
		strings.TrimSpace(cfg.ReasoningExposureOverride) != "" ||
		strings.TrimSpace(cfg.UpstreamProfileOverride) != ""
}

func normalizeCompatToken(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}
