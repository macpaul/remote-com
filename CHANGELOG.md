# Changelog

All notable changes to this project will be documented in this file.

## [v0.1.2] - 2026-05-11

### Added
- **INI Configuration**: Transitioned from JSON to `settings.ini` for better human-readability and manual editing support.
- **UI Enhancements**:
  - Added a "Save Settings" button for manual persistence.
  - Added visual feedback for save operations.
- **Documentation**: Added this `CHANGELOG.md` to track project evolution.

## [v0.1.1] - 2026-05-11

### Added
- **Serial Configuration**: Integrated comprehensive settings for Baud Rate, Data Bits, Parity, Stop Bits, and Flow Control.
- **Transmit Delays**: Added support for character-by-character and line-by-line millisecond delays to support legacy hardware.
- **UI Enhancements**: Improved form layout for easier configuration.
- **Validation**: Implemented strict TCP port range validation (1-65535) on both frontend and backend.
- **Automation**: Added `release.sh` for streamlined cross-platform builds and GitHub asset uploading.
- **Testing**: Added unit tests for port validation and engine management logic.

### Fixed
- **Deadlock**: Resolved a critical race condition in the Manager locking mechanism.
- **Permissions**: Normalized file modes for generated Wails runtime files.

## [v0.1.0] - 2026-05-10

### Added
- **Core Engine**: Initial implementation of the Serial-to-SSH redirection bridge.
- **SSH Server**: Embedded standalone SSH server with password authentication.
- **GUI**: Modern React-based dashboard built with Wails.
- **Persistence**: Basic management layer for TCP-to-Serial port bindings.
- **Cross-Platform**: Support for Windows (AMD64), macOS (ARM64/Intel), and Linux.

[v0.1.1]: https://github.com/macpaul/remote-com/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/macpaul/remote-com/releases/tag/v0.1.0
