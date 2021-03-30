# Changelog
项目更改记录

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