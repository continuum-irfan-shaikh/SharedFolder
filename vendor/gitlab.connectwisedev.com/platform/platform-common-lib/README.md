<p align="center">
<img height=70px src="docs/images/continuum-logo.png">
<img height=70px src="docs/images/Go-Logo_Blue.png">
</p>

# platform-common-lib

This repo is contains all the common components used in the Continuum.

### Repository Contact

- [Common-Frameworks](Common-Frameworks@continuum.net)
- [Continuum Engineering](Project-Juno@continuum.net)

### Common Components
| Component / Module       | Jira/Wiki     | Description                             |
| ------------------------ |:-------------:| :-------------------------------------- |
| [Logger](src/runtime/logger) | [CF-1](https://continuum.atlassian.net/browse/CF-1) | Gray-log friendly Logging framework |
| [Communication-UDP](src/communication/udp) | [CF-18](https://continuum.atlassian.net/browse/CF-18) | UDP Communnication |
| [Communication-HTTP](src/communication/http) | [CF-6](https://continuum.atlassian.net/browse/CF-6) | HTTP Communnication |
| [Binary Version Embedding](src/app) |[Wiki](https://continuum.atlassian.net/wiki/spaces/C2E/pages/1454704686/Continuum+2.0+-+Binary+Version+Embedding) |    Platform agnostic way to get version of binaries by embedding a version information (Version-Info file) along with binary details. |
| [Application Level Metrics](src/metric) | [CF-6](https://continuum.atlassian.net/browse/CF-6) | Collect and Publish Application level metrics |
| [Sanitize](src/sanitize) | [CF-14](https://continuum.atlassian.net/browse/CF-14) | Helper functions to sanitize input strings to protect against XSS, CSV injection, embedded HTML, etc |
| [Validate](src/validate/is) | [CF-14](https://continuum.atlassian.net/browse/CF-14) | Helper functions to validate input strings are either Number, Phone Number, Alpha, Alpha Numeric, UUID etc  |
| [Execute With](src/exec/with) | [CF-7](https://continuum.atlassian.net/browse/CF-7) | Helper functions to execute functions |
| [Runtime Utility](src/runtime/util) | [CF-7](https://continuum.atlassian.net/browse/CF-7) | Helper utility functions to find commonly used values 
| [Crypto Manager](src/cryptomgr) |[Wiki](https://continuum.atlassian.net/wiki/spaces/C2E/pages/946930526/Agent+Core+Sensitive+Data+Encryption) |    Common framework to manage asymmetric cryptography keys. |
| [Downloader](src/downloader) | [CF-36](https://continuum.atlassian.net/browse/CF-36)| Download Manager to download files from Internet using Grab |
| [Redis](src/redis) | [Wiki](https://continuum.atlassian.net/wiki/spaces/EN/pages/1672348148/Device+Down+-+Redis+-+SDD) | Go client for Redis |
| [Utils](src/utils) | [CF-7](https://continuum.atlassian.net/browse/CF-7) | Helper utility functions to perform general operations |
| [Circuit Breaker](src/circuit) | [CF-58](https://continuum.atlassian.net/browse/CF-58) | Circuit breaker wrapper |
| [SQL Db](src/db) | [ITSM-2518](https://continuum.atlassian.net/browse/ITSM-2518) | SQL Db |
| [WAL](src/wal) | [CF-68](https://continuum.atlassian.net/browse/CF-68) | Write Ahead logs |
| [web](src/web) | [CF-67](https://continuum.atlassian.net/browse/CF-67) | Light weight wrapper on top of mux to create server, routes and handle middle-ware functionality |
| [retry](src/retry) | [ZIP-1013](https://continuum.atlassian.net/browse/ZIP-1013) | Wrapper over avast/retry-go framework which provides generic logic for retrying tasks based on different retry strategies |
| [Email Notifications](src/notifications/email) | [IAM-439](https://continuum.atlassian.net/browse/IAM-439) | Email Service Wrapper over AWS SES to send email notifications to users.

### Depricated Packages

| Component / Module            | Depricated On | Removed On| Replacement |
| ------------------------ |:-------------:|:-------------:| :-------------------------------------- |
| [src/logger](src/logger)      | May 30, 2019  | Feb 25, 2020  | [Logger](src/runtime/logger) |
| [src/logging](src/logging)    | May 30, 2019  | Feb 25, 2020  | [Logger](src/runtime/logger) |
| [src/instrumentation](src/instrumentation) | July 30, 2019  | Feb 25, 2020  | [Application Level Metrics](src/metric) |
| [src/kafka](src/kafka) | July 30, 2019  | TBD  | [Infrastructure Messaging](https://gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/tree/master/messaging) |
| [Cherwell](src/cherwell) | October 10, 2019  | TBD  | This will moved to ITSM integration |
| [web/microService](src/web/microService) | May 30, 2019  | TBD  | [web/rest package](src/web/rest) |


### Contiribution

Every one in Continuum can contribute in this repository

#### Contribution Rules

- Package should have clear interfaces and documentation
- Package should have Unit test cases
- Package should have an Example
- Package should have Mocks, if required
- Package should have README.md file
- Package Readme link should be added in the main readme file
- Package should be listed as a Common Component or associated Wiki
- Package should have list of 3rd party libraries along with License : Ensure it is complaint with allowed continuum license
- Please add [Lokesh Jain](mailto:lokesh.jain@continuum.net) and [Nikhil Bhide](mailto:nikhil.bhide@continuum.net) along with others as the reviewer, it will help us to manage common components

Example : [Logger](src/runtime/logger)
