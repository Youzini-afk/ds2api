import { ShieldAlert } from 'lucide-react'
import { useMemo, useState } from 'react'

function normalizeToken(value) {
    return String(value || '').trim().toLowerCase()
}

function resolveCompatPreview(compat) {
    const preset = normalizeToken(compat?.preset) === 'shallowseek_compat' ? 'shallowseek_compat' : 'default'
    const effective = {
        reasoner_prompt_mode: 'default',
        reasoning_exposure: 'always',
        upstream_profile: 'android',
    }

    if (preset === 'shallowseek_compat') {
        effective.reasoner_prompt_mode = 'end_of_thinking'
        effective.reasoning_exposure = 'request_opt_in'
        effective.upstream_profile = 'web'
    }

    const reasonerOverride = normalizeToken(compat?.reasoner_prompt_mode_override)
    if (reasonerOverride === 'default' || reasonerOverride === 'end_of_thinking') {
        effective.reasoner_prompt_mode = reasonerOverride
    }

    const reasoningOverride = normalizeToken(compat?.reasoning_exposure_override)
    if (reasoningOverride === 'always' || reasoningOverride === 'request_opt_in') {
        effective.reasoning_exposure = reasoningOverride
    }

    const upstreamOverride = normalizeToken(compat?.upstream_profile_override)
    if (upstreamOverride === 'android' || upstreamOverride === 'web') {
        effective.upstream_profile = upstreamOverride
    }

    return effective
}

function displayCompatValue(t, field, value) {
    const normalized = normalizeToken(value)
    if (field === 'reasoner_prompt_mode') {
        if (normalized === 'end_of_thinking') return t('settings.compatReasonerModeEndOfThinking')
        return t('settings.compatReasonerModeDefault')
    }
    if (field === 'reasoning_exposure') {
        if (normalized === 'request_opt_in') return t('settings.compatReasoningExposureRequestOptIn')
        return t('settings.compatReasoningExposureAlways')
    }
    if (field === 'upstream_profile') {
        if (normalized === 'web') return t('settings.compatUpstreamProfileWeb')
        return t('settings.compatUpstreamProfileAndroid')
    }
    return value
}

export default function CompatibilitySection({ t, form, setForm }) {
    const [advancedOpen, setAdvancedOpen] = useState(false)
    const compat = form.compat || {}
    const effective = useMemo(() => resolveCompatPreview(compat), [compat])

    const updateCompat = (patch) => {
        setForm((prev) => ({
            ...prev,
            compat: {
                ...prev.compat,
                ...patch,
            },
        }))
    }

    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <div className="flex items-center gap-2">
                <ShieldAlert className="w-4 h-4 text-muted-foreground" />
                <h3 className="font-semibold">{t('settings.compatTitle')}</h3>
            </div>
            <p className="text-sm text-muted-foreground">{t('settings.compatDesc')}</p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.compatPreset')}</span>
                    <select
                        value={compat.preset || 'default'}
                        onChange={(e) => updateCompat({ preset: e.target.value })}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    >
                        <option value="default">{t('settings.compatPresetDefault')}</option>
                        <option value="shallowseek_compat">{t('settings.compatPresetShallowseek')}</option>
                    </select>
                </label>

                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.compatWideInputStrictOutput')}</span>
                    <div className="h-[42px] px-3 rounded-lg border border-border bg-background flex items-center justify-between gap-3">
                        <span className="text-xs text-muted-foreground">{t('settings.compatWideInputStrictOutputDesc')}</span>
                        <input
                            type="checkbox"
                            checked={Boolean(compat.wide_input_strict_output)}
                            onChange={(e) => updateCompat({ wide_input_strict_output: e.target.checked })}
                            className="w-4 h-4"
                        />
                    </div>
                </label>
            </div>

            <div className="flex items-center justify-between gap-4 rounded-lg border border-border bg-background px-3 py-3">
                <label className="text-sm font-medium">{t('settings.stripReferenceMarkers')}</label>
                <button
                    type="button"
                    role="switch"
                    aria-checked={compat.strip_reference_markers ?? true}
                    onClick={() => updateCompat({ strip_reference_markers: !(compat.strip_reference_markers ?? true) })}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                        compat.strip_reference_markers ?? true ? 'bg-primary' : 'bg-muted'
                    }`}
                >
                    <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                            compat.strip_reference_markers ?? true ? 'translate-x-6' : 'translate-x-1'
                        }`}
                    />
                </button>
            </div>

            <div className="rounded-lg border border-border bg-background px-3 py-2 text-xs text-muted-foreground space-y-1">
                <div>
                    {t('settings.compatEffectiveReasonerPromptMode')}: <span className="text-foreground">{displayCompatValue(t, 'reasoner_prompt_mode', effective.reasoner_prompt_mode)}</span>
                </div>
                <div>
                    {t('settings.compatEffectiveReasoningExposure')}: <span className="text-foreground">{displayCompatValue(t, 'reasoning_exposure', effective.reasoning_exposure)}</span>
                </div>
                <div>
                    {t('settings.compatEffectiveUpstreamProfile')}: <span className="text-foreground">{displayCompatValue(t, 'upstream_profile', effective.upstream_profile)}</span>
                </div>
            </div>

            <div className="space-y-3">
                <button
                    type="button"
                    onClick={() => setAdvancedOpen((prev) => !prev)}
                    className="text-xs px-3 py-1.5 rounded-md border border-border hover:bg-muted/50"
                >
                    {advancedOpen ? t('settings.compatHideAdvanced') : t('settings.compatShowAdvanced')}
                </button>
                <p className="text-xs text-muted-foreground">{t('settings.compatAdvancedHint')}</p>
                {advancedOpen && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <label className="text-sm space-y-2">
                            <span className="text-muted-foreground">{t('settings.compatReasonerPromptModeOverride')}</span>
                            <select
                                value={compat.reasoner_prompt_mode_override || ''}
                                onChange={(e) => updateCompat({ reasoner_prompt_mode_override: e.target.value })}
                                className="w-full bg-background border border-border rounded-lg px-3 py-2"
                            >
                                <option value="">{t('settings.compatInheritPreset')}</option>
                                <option value="default">{t('settings.compatReasonerModeDefault')}</option>
                                <option value="end_of_thinking">{t('settings.compatReasonerModeEndOfThinking')}</option>
                            </select>
                        </label>

                        <label className="text-sm space-y-2">
                            <span className="text-muted-foreground">{t('settings.compatReasoningExposureOverride')}</span>
                            <select
                                value={compat.reasoning_exposure_override || ''}
                                onChange={(e) => updateCompat({ reasoning_exposure_override: e.target.value })}
                                className="w-full bg-background border border-border rounded-lg px-3 py-2"
                            >
                                <option value="">{t('settings.compatInheritPreset')}</option>
                                <option value="always">{t('settings.compatReasoningExposureAlways')}</option>
                                <option value="request_opt_in">{t('settings.compatReasoningExposureRequestOptIn')}</option>
                            </select>
                        </label>

                        <label className="text-sm space-y-2">
                            <span className="text-muted-foreground">{t('settings.compatUpstreamProfileOverride')}</span>
                            <select
                                value={compat.upstream_profile_override || ''}
                                onChange={(e) => updateCompat({ upstream_profile_override: e.target.value })}
                                className="w-full bg-background border border-border rounded-lg px-3 py-2"
                            >
                                <option value="">{t('settings.compatInheritPreset')}</option>
                                <option value="android">{t('settings.compatUpstreamProfileAndroid')}</option>
                                <option value="web">{t('settings.compatUpstreamProfileWeb')}</option>
                            </select>
                        </label>
                    </div>
                )}
            </div>
        </div>
    )
}
