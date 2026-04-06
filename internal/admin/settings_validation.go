package admin

import (
	"fmt"
	"strings"

	"ds2api/internal/config"
)

func normalizeSettingsConfig(c *config.Config) {
	if c == nil {
		return
	}
	c.Admin.PasswordHash = strings.TrimSpace(c.Admin.PasswordHash)
	c.Embeddings.Provider = strings.TrimSpace(c.Embeddings.Provider)
	c.Compat.Preset = strings.ToLower(strings.TrimSpace(c.Compat.Preset))
	c.Compat.ReasonerPromptModeOverride = strings.ToLower(strings.TrimSpace(c.Compat.ReasonerPromptModeOverride))
	c.Compat.ReasoningExposureOverride = strings.ToLower(strings.TrimSpace(c.Compat.ReasoningExposureOverride))
	c.Compat.UpstreamProfileOverride = strings.ToLower(strings.TrimSpace(c.Compat.UpstreamProfileOverride))
}

func validateSettingsConfig(c config.Config) error {
	if err := config.ValidateConfig(c); err != nil {
		return err
	}
	return validateCompatSettings(c.Compat)
}

func validateCompatSettings(compat config.CompatConfig) error {
	if compat.Preset != "" && !config.IsValidCompatPreset(compat.Preset) {
		return fmt.Errorf("compat.preset must be one of: default, shallowseek_compat")
	}
	if compat.ReasonerPromptModeOverride != "" && !config.IsValidCompatReasonerPromptMode(compat.ReasonerPromptModeOverride) {
		return fmt.Errorf("compat.reasoner_prompt_mode_override must be one of: default, end_of_thinking, or empty")
	}
	if compat.ReasoningExposureOverride != "" && !config.IsValidCompatReasoningExposure(compat.ReasoningExposureOverride) {
		return fmt.Errorf("compat.reasoning_exposure_override must be one of: always, request_opt_in, or empty")
	}
	if compat.UpstreamProfileOverride != "" && !config.IsValidCompatUpstreamProfile(compat.UpstreamProfileOverride) {
		return fmt.Errorf("compat.upstream_profile_override must be one of: android, web, or empty")
	}
	return nil
}

func validateRuntimeSettings(runtime config.RuntimeConfig) error {
	return config.ValidateRuntimeConfig(runtime)
}
