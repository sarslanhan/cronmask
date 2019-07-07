[![Build Status](https://travis-ci.org/sarslanhan/cronmask.svg?branch=master)](https://travis-ci.org/sarslanhan/cronmask)
[![codecov](https://codecov.io/gh/sarslanhan/cronmask/branch/master/graph/badge.svg)](https://codecov.io/gh/sarslanhan/cronmask)
[![Go Report Card](https://goreportcard.com/badge/zalando/skipper)](https://goreportcard.com/report/sarslanhan/cronmask)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/zalando/skipper?status.svg)](https://godoc.org/github.com/sarslanhan/cronmask)


# cronmask
Library that supports cron-like expressions to check if a time fits into a time range

For [CRON expressions](https://en.wikipedia.org/wiki/Cron#CRON_expression):

Expressions are expected to be in the same time zone as the system that generates the `time.Time` instances.

You can check the tests for what is possible.

**Unsupported features:**

- [Non-standard characters](https://en.wikipedia.org/wiki/Cron#Non-standard_characters)
- Year field
- Command section
- Text representation of the fields "month" and "day of week"
