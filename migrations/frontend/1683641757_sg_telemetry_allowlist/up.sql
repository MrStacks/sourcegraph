-- This migration was generated by the command `sg telemetry add`
INSERT INTO event_logs_export_allowlist (event_name) VALUES (UNNEST('{CodySignup,VSCodeInstall,VSCodeMarketplace,TryCodyWeb,TryCodyWebOnboardingDisplayed}'::TEXT[])) ON CONFLICT DO NOTHING;