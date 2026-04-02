import { useState } from 'react'

export default function CompatibilitySection({ t, form, setForm }) {
    const [advancedOpen, setAdvancedOpen] = useState(false)
    const compat = form.compat || {}

    const updateCompat = (patch) => {
        setForm(prev => ({
            ...prev,
            compat: {
                ...prev.compat,
                ...patch,
            },
        }))
    }

    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <div className="space-y-1">
                <h3 className="font-semibold">{t('settings.compatTitle')}</h3>
                <p className="text-xs text-muted-foreground">{t('settings.compatDesc')}</p>
            </div>

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
                    <div className="h-[42px] px-3 rounded-lg border border-border bg-background flex items-center justify-between">
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

            <div className="rounded-lg border border-border bg-background px-3 py-2 text-xs text-muted-foreground space-y-1">
                <div>
                    {t('settings.compatEffectiveReasonerPromptMode')}: <span className="font-mono text-foreground">{compat.effective_reasoner_prompt_mode || 'default'}</span>
                </div>
                <div>
                    {t('settings.compatEffectiveReasoningExposure')}: <span className="font-mono text-foreground">{compat.effective_reasoning_exposure || 'always'}</span>
                </div>
                <div>
                    {t('settings.compatEffectiveUpstreamProfile')}: <span className="font-mono text-foreground">{compat.effective_upstream_profile || 'android'}</span>
                </div>
            </div>

            <div className="space-y-3">
                <button
                    type="button"
                    onClick={() => setAdvancedOpen(prev => !prev)}
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
                                <option value="default">default</option>
                                <option value="end_of_thinking">end_of_thinking</option>
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
                                <option value="always">always</option>
                                <option value="request_opt_in">request_opt_in</option>
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
                                <option value="android">android</option>
                                <option value="web">web</option>
                            </select>
                        </label>
                    </div>
                )}
            </div>
        </div>
    )
}
