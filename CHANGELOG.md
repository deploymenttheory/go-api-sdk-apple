# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.10.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.9.0...v0.10.0) (2026-09-01)


### Features

* **notary:** add S3 upload and end-to-end SubmitAndWait workflow ([89d1a00](https://github.com/deploymenttheory/go-sdk-appleservices/commit/89d1a0036978e883c67ef1812df682a79fa4910a))
* **notary:** add S3 upload and end-to-end SubmitAndWait workflow ([5fe5804](https://github.com/deploymenttheory/go-sdk-appleservices/commit/5fe580468426675def6a22a3f96f073580e80ea7))

## [Unreleased]

### Features

* **notary:** add S3 upload and end-to-end `SubmitAndWait` workflow. The notary service can now notarize a file from hash to verdict in one call: it submits the software, uploads it to the S3 bucket Apple hands out (AWS Signature Version 4, standard-library only, no AWS SDK dependency), polls for the verdict, and fetches and parses the developer log on a rejection. A new `notary/upload` package holds the stand-alone SigV4 uploader, and `examples/notary/` shows the flow.

## [0.9.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.8.1...v0.9.0) (2026-08-30)


### Features

* **axm:** add device management service migration and device release ([a06fba4](https://github.com/deploymenttheory/go-sdk-appleservices/commit/a06fba43d45864a081b3cde457d1ca1445d10eee))
* **axm:** add device management service migration and device release ([de438ba](https://github.com/deploymenttheory/go-sdk-appleservices/commit/de438ba8355e0348f2be83ea3ceb124e7113c969))

## [0.8.1](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.8.0...v0.8.1) (2026-08-29)


### Bug Fixes

* **module:** declare the module at the repository's import path ([31c6d0d](https://github.com/deploymenttheory/go-sdk-appleservices/commit/31c6d0d94bb5ac1d0a4af756fddd0739e6b25c38))
* **module:** declare the module at the repository's import path ([71db06a](https://github.com/deploymenttheory/go-sdk-appleservices/commit/71db06a42269ad971a7ff06032cda7a6c114245b))

## [0.8.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.7.0...v0.8.0) (2026-08-03)


### Features

* added abm v2.2 updates - https://developer.apple.com/documentation/apple-school-and-business-manager-api/apple-school-manager-and-apple-business-api-changelog ([e552abd](https://github.com/deploymenttheory/go-sdk-appleservices/commit/e552abd1310d96bf6d49047239f7cd2abd09efcb))
* added abm v2.2 updates - https://developer.apple.com/documentation/apple-school-and-business-manager-api/apple-school-manager-and-apple-business-api-changelog ([07c6153](https://github.com/deploymenttheory/go-sdk-appleservices/commit/07c6153aac385420b5654474a258c157ae153efa))

## [0.7.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.6.0...v0.7.0) (2026-07-17)


### Features

* **device_management:** add generated Apple MDM/DDM SDK built from apple/device-management specs ([c6cda7b](https://github.com/deploymenttheory/go-sdk-appleservices/commit/c6cda7b261c7d785ca7746732da1af8741984dd4))
* **device_management:** add generated Apple MDM/DDM SDK built from apple/device-management specs ([8772aa7](https://github.com/deploymenttheory/go-sdk-appleservices/commit/8772aa70e33b233591ffdc1f524576e065d11f28))

## [0.6.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.5.0...v0.6.0) (2026-06-07)


### Features

* added abm v2.1 updates - ([5f336b8](https://github.com/deploymenttheory/go-sdk-appleservices/commit/5f336b8a7052d90e79e18ab46c0fe5f1ba1ee1b1))
* added abm v2.1 updates - https://developer.apple.com/documentation/apple-school-and-business-manager-api/apple-school-manager-and-apple-business-api-changelog ([6990259](https://github.com/deploymenttheory/go-sdk-appleservices/commit/69902591ecf93ad48b6135d80e89e831b8d0c6f6))
* added support for notary api ([7d226e3](https://github.com/deploymenttheory/go-sdk-appleservices/commit/7d226e34ced539474bc6b8d900100ed4ab70acdb))
* added support for notary api ([f309276](https://github.com/deploymenttheory/go-sdk-appleservices/commit/f309276db85d98288c78ca90a2c9ffe7fcfe96c2))

## [0.5.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.4.0...v0.5.0) (2026-04-27)


### Features

* enhance Apple Business Manager SDK with new API endpoints and s… ([f4c9d8f](https://github.com/deploymenttheory/go-sdk-appleservices/commit/f4c9d8f256f8fada24d614dcab7a612e5e96da0d))
* enhance Apple Business Manager SDK with new API endpoints and services ([c4cd610](https://github.com/deploymenttheory/go-sdk-appleservices/commit/c4cd6108022da47bfe40c029906df5e24eb24ce1))

## [0.4.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.3.1...v0.4.0) (2026-04-21)


### Features

* add apple_update_cdn SDK for firmware discovery and IPSW downloads ([6e02b88](https://github.com/deploymenttheory/go-sdk-appleservices/commit/6e02b88b5400c8ca729507c54abdc8103746133e))


### Bug Fixes

* update device attribute constants and field names ([7ef7c7a](https://github.com/deploymenttheory/go-sdk-appleservices/commit/7ef7c7a12c67f32aef21e8b483047c03673131d9))
* update device attribute constants and field names ([23d8521](https://github.com/deploymenttheory/go-sdk-appleservices/commit/23d8521902126bfe7f9c6d04de349814de683b3e))
* update README for API client usage ([87d759b](https://github.com/deploymenttheory/go-sdk-appleservices/commit/87d759b0adad26c24dd1f330189dfa3cf79d000b))
* update README for API client usage ([468a261](https://github.com/deploymenttheory/go-sdk-appleservices/commit/468a261c221dfe6a6110fc252da3539a7315c136))

## [0.3.1](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.3.0...v0.3.1) (2026-03-29)


### Bug Fixes

* Apple Business Manager client methods to support optional configuration options ([d60be52](https://github.com/deploymenttheory/go-sdk-appleservices/commit/d60be52417b5d73c449a37a9b2c496f75257fab2))
* pagination methods in client and services ([7576d4f](https://github.com/deploymenttheory/go-sdk-appleservices/commit/7576d4f9a3143dabb0c3c1ca00284030b73cb0e5))

## [0.3.1](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.3.0...v0.3.1) (2026-03-27)


### Bug Fixes

* Apple Business Manager client methods to support optional configuration options ([d60be52](https://github.com/deploymenttheory/go-sdk-appleservices/commit/d60be52417b5d73c449a37a9b2c496f75257fab2))
* pagination methods in client and services ([7576d4f](https://github.com/deploymenttheory/go-sdk-appleservices/commit/7576d4f9a3143dabb0c3c1ca00284030b73cb0e5))

## [0.3.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.2.0...v0.3.0) (2025-11-09)


### Features

* add microsoft_mac_apps_version_tracker ([0ea6555](https://github.com/deploymenttheory/go-sdk-appleservices/commit/0ea65555a7012c48901a4d10466f7e6f9bcd1eee))
* added func Get AppleCare Coverage Information for a Device with tests and examples ([0f0d252](https://github.com/deploymenttheory/go-sdk-appleservices/commit/0f0d2521f48d053960ee5e54f064d038297b37e4))

## [0.2.0](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.1.1...v0.2.0) (2025-10-22)


### Features

* added v3 for testing ([1033b91](https://github.com/deploymenttheory/go-sdk-appleservices/commit/1033b91a5272508a7333794ab5d07280c6611b51))


### Bug Fixes

* examples for axm ([d440074](https://github.com/deploymenttheory/go-sdk-appleservices/commit/d440074f23fe32a6bd836ae5c93eb36af0463ba6))
* examples tweaks ([da6a634](https://github.com/deploymenttheory/go-sdk-appleservices/commit/da6a634d0f785ae0f51e5982160c0129322a132b))
* for constants usage ([68b48e4](https://github.com/deploymenttheory/go-sdk-appleservices/commit/68b48e4629575723c051c48988a09fffc0dd0a70))
* more struct dupes ([e6c73b8](https://github.com/deploymenttheory/go-sdk-appleservices/commit/e6c73b8c04b6d02415dd79a93e15e5f8f56955ae))
* query struct dupe ([34cfa2f](https://github.com/deploymenttheory/go-sdk-appleservices/commit/34cfa2f9a62df538c2f16da70c1efb20d92befd1))
* removed reudundant oauth auth code ([4f0274f](https://github.com/deploymenttheory/go-sdk-appleservices/commit/4f0274f8175fadf03afc98dbfa4984926ac17b2e))
* restructure ([3cb4069](https://github.com/deploymenttheory/go-sdk-appleservices/commit/3cb40692b672787f77cc0bb676dbb90eb383170b))
* struct reorg ([da77d79](https://github.com/deploymenttheory/go-sdk-appleservices/commit/da77d79e7af59dfb3ec895582f2d3a7ebde09431))

## [0.1.1](https://github.com/deploymenttheory/go-sdk-appleservices/compare/v0.1.0...v0.1.1) (2025-10-09)


### Bug Fixes

* rl pat token var ([ad16e08](https://github.com/deploymenttheory/go-sdk-appleservices/commit/ad16e08f8f03a3e9ecc1a659d388498f14bc4a69))
* tidy up ([9e878fc](https://github.com/deploymenttheory/go-sdk-appleservices/commit/9e878fc18c9bd8f8c45ecc50f1167d4652627e33))

## [Unreleased]

### Added

- Added xyz [@your_username](https://github.com/your_username)

### Fixed

- Fixed zyx [@your_username](https://github.com/your_username)

## [1.1.0] - 2021-06-23

### Added

- Added x [@your_username](https://github.com/your_username)

### Changed

- Changed y [@your_username](https://github.com/your_username)

## [1.0.0] - 2021-06-20

### Added

- Inititated y [@your_username](https://github.com/your_username)
- Inititated z [@your_username](https://github.com/your_username)
