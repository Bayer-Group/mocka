# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Revive for linting Go code
- Dlv for debugging Go code
- Additional nil checks to prevent internal panics
- Mutex to sandbox
- RWMutex to mockFunction

### Changed
- Examples to have more use cases
- Sandbox to use pointer to stubs to not copy mutex value

### Fixed
- Equality checks in WithArgs to use reflect.DeepEqual to compare slices

## [1.0.0]
### Added
- Initial public release of this package
