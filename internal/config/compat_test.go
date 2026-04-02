package config

import "testing"

func TestResolveCompatDefaultPreset(t *testing.T) {
	eff := ResolveCompat(CompatConfig{})
	if eff.Preset != CompatPresetDefault {
		t.Fatalf("expected preset=%q, got %q", CompatPresetDefault, eff.Preset)
	}
	if eff.ReasonerPromptMode != CompatReasonerPromptDefault {
		t.Fatalf("expected reasoner_prompt_mode=%q, got %q", CompatReasonerPromptDefault, eff.ReasonerPromptMode)
	}
	if eff.ReasoningExposure != CompatReasoningAlways {
		t.Fatalf("expected reasoning_exposure=%q, got %q", CompatReasoningAlways, eff.ReasoningExposure)
	}
	if eff.UpstreamProfile != CompatUpstreamAndroid {
		t.Fatalf("expected upstream_profile=%q, got %q", CompatUpstreamAndroid, eff.UpstreamProfile)
	}
}

func TestResolveCompatShallowseekPreset(t *testing.T) {
	eff := ResolveCompat(CompatConfig{Preset: CompatPresetShallowseek})
	if eff.ReasonerPromptMode != CompatReasonerPromptEndThink {
		t.Fatalf("expected reasoner_prompt_mode=%q, got %q", CompatReasonerPromptEndThink, eff.ReasonerPromptMode)
	}
	if eff.ReasoningExposure != CompatReasoningRequestOptIn {
		t.Fatalf("expected reasoning_exposure=%q, got %q", CompatReasoningRequestOptIn, eff.ReasoningExposure)
	}
	if eff.UpstreamProfile != CompatUpstreamWeb {
		t.Fatalf("expected upstream_profile=%q, got %q", CompatUpstreamWeb, eff.UpstreamProfile)
	}
}

func TestResolveCompatOverridePriority(t *testing.T) {
	eff := ResolveCompat(CompatConfig{
		Preset:                     CompatPresetShallowseek,
		ReasonerPromptModeOverride: CompatReasonerPromptDefault,
		ReasoningExposureOverride:  CompatReasoningAlways,
		UpstreamProfileOverride:    CompatUpstreamAndroid,
	})
	if eff.ReasonerPromptMode != CompatReasonerPromptDefault {
		t.Fatalf("expected reasoner_prompt_mode override to win, got %q", eff.ReasonerPromptMode)
	}
	if eff.ReasoningExposure != CompatReasoningAlways {
		t.Fatalf("expected reasoning_exposure override to win, got %q", eff.ReasoningExposure)
	}
	if eff.UpstreamProfile != CompatUpstreamAndroid {
		t.Fatalf("expected upstream_profile override to win, got %q", eff.UpstreamProfile)
	}
}

func TestResolveCompatIgnoresInvalidOverrides(t *testing.T) {
	eff := ResolveCompat(CompatConfig{
		Preset:                     CompatPresetShallowseek,
		ReasonerPromptModeOverride: "invalid",
		ReasoningExposureOverride:  "invalid",
		UpstreamProfileOverride:    "invalid",
	})
	if eff.ReasonerPromptMode != CompatReasonerPromptEndThink {
		t.Fatalf("expected preset fallback reasoner prompt mode, got %q", eff.ReasonerPromptMode)
	}
	if eff.ReasoningExposure != CompatReasoningRequestOptIn {
		t.Fatalf("expected preset fallback reasoning exposure, got %q", eff.ReasoningExposure)
	}
	if eff.UpstreamProfile != CompatUpstreamWeb {
		t.Fatalf("expected preset fallback upstream profile, got %q", eff.UpstreamProfile)
	}
}
