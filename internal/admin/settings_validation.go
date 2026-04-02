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
	if c.Admin.JWTExpireHours != 0 && (c.Admin.JWTExpireHours < 1 || c.Admin.JWTExpireHours > 720) {
		return fmt.Errorf("admin.jwt_expire_hours must be between 1 and 720")
	}
	if err := validateRuntimeSettings(c.Runtime); err != nil {
		return err
	}
	if c.Responses.StoreTTLSeconds != 0 && (c.Responses.StoreTTLSeconds < 30 || c.Responses.StoreTTLSeconds > 86400) {
		return fmt.Errorf("responses.store_ttl_seconds must be between 30 and 86400")
	}
	if c.Embeddings.Provider != "" && strings.TrimSpace(c.Embeddings.Provider) == "" {
		return fmt.Errorf("embeddings.provider cannot be empty")
	}
	if err := validateCompatSettings(c.Compat); err != nil {
		return err
	}
	return nil
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
	if runtime.AccountMaxInflight != 0 && (runtime.AccountMaxInflight < 1 || runtime.AccountMaxInflight > 256) {
		return fmt.Errorf("runtime.account_max_inflight must be between 1 and 256")
	}
	if runtime.AccountMaxQueue != 0 && (runtime.AccountMaxQueue < 1 || runtime.AccountMaxQueue > 200000) {
		return fmt.Errorf("runtime.account_max_queue must be between 1 and 200000")
	}
	if runtime.GlobalMaxInflight != 0 && (runtime.GlobalMaxInflight < 1 || runtime.GlobalMaxInflight > 200000) {
		return fmt.Errorf("runtime.global_max_inflight must be between 1 and 200000")
	}
	if runtime.TokenRefreshIntervalHours != 0 && (runtime.TokenRefreshIntervalHours < 1 || runtime.TokenRefreshIntervalHours > 720) {
		return fmt.Errorf("runtime.token_refresh_interval_hours must be between 1 and 720")
	}
	if runtime.AccountMaxInflight > 0 && runtime.GlobalMaxInflight > 0 && runtime.GlobalMaxInflight < runtime.AccountMaxInflight {
		return fmt.Errorf("runtime.global_max_inflight must be >= runtime.account_max_inflight")
	}
	return nil
}
