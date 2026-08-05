# Changelog

## 2.1.0 — 2026-08-04

### Added

- Legacy embedding alias `caedral-embed` alongside canonical `caedral-embed-e1-small-v1`
- `InputType` and `EncodingFormat` fields on `EmbeddingCreateRequest` for OpenRouter provider parity

## 2.0.0 — 2026-08-04

### Breaking

- Embeddings default to `caedral-embed-e1-small-v1` with native 384 dimensions (E1 Small)
- Client-side validation for supported embedding models and dimensions

## 1.0.0 — 2026-07-29

First stable release of the official Caedral Go SDK.

- Module path: `github.com/trycaedral/caedral-go`
- Install: `go get github.com/trycaedral/caedral-go@v1.0.0`
