package admin

import (
	"fmt"
	"strings"

	"ds2api/internal/config"
)

func boolFrom(v any) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return strings.ToLower(strings.TrimSpace(x)) == "true"
	default:
		return false
	}
}

func normalizeOptionalToken(v any) string {
	if v == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
}

type settingsCompatPatch struct {
	HasWideInputStrictOutput      bool
	WideInputStrictOutput         bool
	HasStripReferenceMarkers      bool
	StripReferenceMarkers         bool
	HasPreset                     bool
	Preset                        string
	HasReasonerPromptModeOverride bool
	ReasonerPromptModeOverride    string
	HasReasoningExposureOverride  bool
	ReasoningExposureOverride     string
	HasUpstreamProfileOverride    bool
	UpstreamProfileOverride       string
}

func parseSettingsUpdateRequest(req map[string]any) (*config.AdminConfig, *config.RuntimeConfig, *config.ResponsesConfig, *config.EmbeddingsConfig, *config.AutoDeleteConfig, *settingsCompatPatch, map[string]string, map[string]string, error) {
	var (
		adminCfg      *config.AdminConfig
		runtimeCfg    *config.RuntimeConfig
		respCfg       *config.ResponsesConfig
		embCfg        *config.EmbeddingsConfig
		autoDeleteCfg *config.AutoDeleteConfig
		compatCfg     *settingsCompatPatch
		claudeMap     map[string]string
		aliasMap      map[string]string
	)

	if raw, ok := req["admin"].(map[string]any); ok {
		cfg := &config.AdminConfig{}
		if v, exists := raw["jwt_expire_hours"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("admin.jwt_expire_hours", n, 1, 720, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.JWTExpireHours = n
		}
		adminCfg = cfg
	}

	if raw, ok := req["runtime"].(map[string]any); ok {
		cfg := &config.RuntimeConfig{}
		if v, exists := raw["account_max_inflight"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("runtime.account_max_inflight", n, 1, 256, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.AccountMaxInflight = n
		}
		if v, exists := raw["account_max_queue"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("runtime.account_max_queue", n, 1, 200000, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.AccountMaxQueue = n
		}
		if v, exists := raw["global_max_inflight"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("runtime.global_max_inflight", n, 1, 200000, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.GlobalMaxInflight = n
		}
		if v, exists := raw["token_refresh_interval_hours"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("runtime.token_refresh_interval_hours", n, 1, 720, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.TokenRefreshIntervalHours = n
		}
		if cfg.AccountMaxInflight > 0 && cfg.GlobalMaxInflight > 0 && cfg.GlobalMaxInflight < cfg.AccountMaxInflight {
			return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("runtime.global_max_inflight must be >= runtime.account_max_inflight")
		}
		runtimeCfg = cfg
	}

	if raw, ok := req["compat"].(map[string]any); ok {
		cfg := &settingsCompatPatch{}
		if v, exists := raw["wide_input_strict_output"]; exists {
			cfg.HasWideInputStrictOutput = true
			cfg.WideInputStrictOutput = boolFrom(v)
		}
		if v, exists := raw["strip_reference_markers"]; exists {
			cfg.HasStripReferenceMarkers = true
			cfg.StripReferenceMarkers = boolFrom(v)
		}
		if v, exists := raw["preset"]; exists {
			preset := normalizeOptionalToken(v)
			if preset != "" && !config.IsValidCompatPreset(preset) {
				return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("compat.preset must be one of: default, shallowseek_compat")
			}
			cfg.HasPreset = true
			cfg.Preset = preset
		}
		if v, exists := raw["reasoner_prompt_mode_override"]; exists {
			mode := normalizeOptionalToken(v)
			if mode != "" && !config.IsValidCompatReasonerPromptMode(mode) {
				return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("compat.reasoner_prompt_mode_override must be one of: default, end_of_thinking, or empty")
			}
			cfg.HasReasonerPromptModeOverride = true
			cfg.ReasonerPromptModeOverride = mode
		}
		if v, exists := raw["reasoning_exposure_override"]; exists {
			mode := normalizeOptionalToken(v)
			if mode != "" && !config.IsValidCompatReasoningExposure(mode) {
				return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("compat.reasoning_exposure_override must be one of: always, request_opt_in, or empty")
			}
			cfg.HasReasoningExposureOverride = true
			cfg.ReasoningExposureOverride = mode
		}
		if v, exists := raw["upstream_profile_override"]; exists {
			profile := normalizeOptionalToken(v)
			if profile != "" && !config.IsValidCompatUpstreamProfile(profile) {
				return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("compat.upstream_profile_override must be one of: android, web, or empty")
			}
			cfg.HasUpstreamProfileOverride = true
			cfg.UpstreamProfileOverride = profile
		}
		compatCfg = cfg
	}

	if raw, ok := req["responses"].(map[string]any); ok {
		cfg := &config.ResponsesConfig{}
		if v, exists := raw["store_ttl_seconds"]; exists {
			n := intFrom(v)
			if err := config.ValidateIntRange("responses.store_ttl_seconds", n, 30, 86400, true); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.StoreTTLSeconds = n
		}
		respCfg = cfg
	}

	if raw, ok := req["embeddings"].(map[string]any); ok {
		cfg := &config.EmbeddingsConfig{}
		if v, exists := raw["provider"]; exists {
			p := strings.TrimSpace(fmt.Sprintf("%v", v))
			if err := config.ValidateTrimmedString("embeddings.provider", p, false); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			cfg.Provider = p
		}
		embCfg = cfg
	}

	if raw, ok := req["claude_mapping"].(map[string]any); ok {
		claudeMap = map[string]string{}
		for k, v := range raw {
			key := strings.TrimSpace(k)
			val := strings.TrimSpace(fmt.Sprintf("%v", v))
			if key == "" || val == "" {
				continue
			}
			claudeMap[key] = val
		}
	}

	if raw, ok := req["model_aliases"].(map[string]any); ok {
		aliasMap = map[string]string{}
		for k, v := range raw {
			key := strings.TrimSpace(k)
			val := strings.TrimSpace(fmt.Sprintf("%v", v))
			if key == "" || val == "" {
				continue
			}
			aliasMap[key] = val
		}
	}

	if raw, ok := req["auto_delete"].(map[string]any); ok {
		cfg := &config.AutoDeleteConfig{}
		if v, exists := raw["mode"]; exists {
			mode := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", v)))
			if err := config.ValidateAutoDeleteMode(mode); err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			if mode == "" {
				mode = "none"
			}
			cfg.Mode = mode
		}
		if v, exists := raw["sessions"]; exists {
			cfg.Sessions = boolFrom(v)
		}
		autoDeleteCfg = cfg
	}

	return adminCfg, runtimeCfg, respCfg, embCfg, autoDeleteCfg, compatCfg, claudeMap, aliasMap, nil
}
