[![Build Status](https://travis-ci.com/vadim-hleif/godiffer.svg?branch=master)](https://travis-ci.com/vadim-hleif/godiffer)
[![Go Report Card](https://goreportcard.com/badge/github.com/vadim-hleif/godiffer)](https://goreportcard.com/report/github.com/vadim-hleif/godiffer)
[![codecov](https://codecov.io/gh/vadim-hleif/godiffer/branch/master/graph/badge.svg)](https://codecov.io/gh/vadim-hleif/godiffer)
[![codebeat badge](https://codebeat.co/badges/8aab8f2d-7652-46f0-a244-9f80dbed64b2)](https://codebeat.co/projects/github-com-vadim-hleif-godiffer-master)

***
- [X] all path and verb combinations in the old specification are present in the new one
- [X] no request parameters are required in the new specification that were not required in the old one
- [X] all required request parameters in the old specification are present in the new one
- [X] all request parameters in the old specification have the same type in the new one
- [X] all request models via $ref don't have any differences
- [X] all response attributes in the old specification are present in the new one
- [X] all response attributes in the old specification have the same type in the new one
- [X] all response models via $ref don't have any differences
- [X] enums validation
- [X] arrays type support
- [X] recursive refs
- [ ] extensions validation


TODO MINOR THINGS
- [ ] common / global parameters (in root of some path / specification)
- [ ] empty value parameters
- [X] ref to not definitions node
- [X] recursive objects
- [ ] reusable enums
- [ ] array of anonymous objects
- [ ] definition in V1 with object in V2 (both cases)
- [ ] reusing responses https://swagger.io/docs/specification/2-0/describing-responses/