项目更改记录
# Changelog
## [0.12.7]
### Changed
- frontlock middleware

## [0.12.6]
### Changed
- remove setting Authorization header in rest_client when having token parameter
- remove HS256 decode jwt token from jwt middleware

## [0.12.5]
### Added
- support []byte type payload in restclient

## [0.12.4]
### Changed
- add contributors to swag

## [0.12.3]
### Added
- cors middleware

## [0.12.2]
### Changed
- add SC_SKIP_TOKEN_CACHE to skip cache in GetAppToken and GetUserToken

## [0.12.1]
### Changed
- remove secret parameter from DefaultJwtMiddleware in jwt middleware

## [0.12.0]
### Changed
- jwt middleware
- auth middleware
- accesscontrol middleware
- restclient
- GetAppToken and GetUserToken in authorization

## [0.11.3]
### Changed
- add token key for restclient kwargs parameter

## [0.11.2]
### Changed
- update swag package

## [0.11.1]
### Changed
- fix the bug in auth middleware

## [0.11.0]
### Added
- frontlock middleware

## [0.10.1]
### Added
- Custom option for cache
### Changed
- add ping mothod when init cache

## [0.10.0]
### Added
- Cache
### Changed
- add redis cache when GetAppToken and GetUserToken in authorization

## [0.9.0]
### Added
- RabbitMQ consumer worker

## [0.8.1]
### Changed
- Modify the svc and api key of Distributed lock client

## [0.8.0]
### Added
- Distributed lock client

## [0.7.2]
### Changed
- Modify the send lanyue message function

## [0.7.1]
### Added
- Save more jwt token claims information to middleware

## [0.7.0]
### Added
- Send lanyue message

## [0.6.0]
### Added
- Send dingtalk message function
- Send rabbitMQ message function
- Get authorization token function
- Request headers can be passed in restclient
### Changed
- restclient options parameter type change

## [0.5.1]
### Fixed
- Fix swag register failed bug

## [0.5.0]
### Added
- swag: register swagger
- auth: add CheckUser function