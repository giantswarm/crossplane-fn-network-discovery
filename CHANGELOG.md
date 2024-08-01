# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2024-08-01

### This version contains breaking changes to the API

As part of this release, the output API has been restructured to allow for
broader extension in the future.

Most `map[string]string` interfaces are now replaced with `map[string]struct`
interfaces meaning the expected `id` will be one level lower on the `id` field.

This change has been carried out to allow more information to be recovered about
the network components such as amazon resource names and other information which
may be useful for understanding how the network looks as part of a composition.

## [0.2.4] - 2024-06-20

- handle service endpoints locally

## [0.2.3] - 2024-06-17

- Fix missing region and provider config

## [0.2.2] - 2024-06-17

- Remove primary cidr block from list of additional cidrs to prevent duplication
- Fix missing region and provider config

## [0.2.1] - 2024-06-17

- Remove empty entries from list

## [0.2.0] - 2024-06-15

- Fix minor issues with grouping when no groups specified

## [0.1.0] - 2024-06-01

- Initial release of function-network-discovery

[Unreleased]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.2.4...v0.3.0
[0.2.4]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.2.1...v0.2.3
[0.2.2]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/crossplane-fn-network-discovery/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/crossplane-fn-network-discovery/releases/tag/v0.2.0
[0.1.0]: https://github.com/giantswarm/crossplane-fn-network-discovery/releases/tag/v0.1.0
