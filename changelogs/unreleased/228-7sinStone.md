# Add universeDomain parameter for GCS client

**Type:** Feature

**Description:**
Introduce a new `universeDomain` key in the BackupStorageLocation `spec.config` map,
allowing users to override the default Google Storage domain (defaults to `"googleapis.com"`).

**Motivation:**
- Support custom or private GCS-compatible endpoints (e.g. testing, on‑prem).
- Align with Helm chart support for `universeDomain`.

**Changes:**
- Updated `object_storage.go` to inject `option.WithUniverseDomain(...)`.
