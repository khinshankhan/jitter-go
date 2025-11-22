# Changelog

All notable changes to this project will be documented in this file.

## [v2.0.1] - 2025-11-18 - 5264b31

**Hotfix** for the v2 module release, fixed borked path.

### Fixed

- Corrected `go.mod` module path for `github.com/khinshankhan/jitter-go/v2`.
- Updated installation instructions in README.

## [v2.0.0] - 2025-11-18 - d51f9d9

**Major release. Breaking changes.**

This version introduces explicit error handling and a more robust configuration validation API.

### Breaking

- `New(Config)` now returns `(Strategy, error)` instead of panicking.
- Existing v1 code must now check the returned error.

### Added

- `ConfigError` type with an `Issues []string` field for inspecting invalid configuration.
- Exposed individual strategy validation errors.
- New public error surface for each strategy implementation.
- Expanded README with examples and usage guidance.

---

## [v1.0.0] - 2025-11-16 - 6c65af9

Initial stable release.

### Added

- Base strategy abstractions.
- Full jitter implementation (default).
- Equal jitter.
- Decorrelated jitter (stateful).
- Naive exponential (no jitter).
- Cap-aware exponential behavior.

---

## [Unreleased]

- (Nothing yet)
