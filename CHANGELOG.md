# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.0.0] - 2024-01-01

### Added
- Initial release of the must_go package
- Basic error handling with `Must(err)` function
- Custom error messages with `MustWithMessage(err, message)`
- HTTP-specific error handling with `MustHTTP(err, statusCode, message)`
- Automatic HTTP error detection with `MustHTTPWithDefault(err)`
- Recovery middleware for HTTP handlers
- Common HTTP error helper functions (404, 400, 401, etc.)
- Generic parsing functions with automatic error handling
- Comprehensive test suite
- Example application demonstrating all features
- Complete documentation and README 