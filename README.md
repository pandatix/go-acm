<div align="center">
	<h1>Go-ACM</h1>
	<a href="https://pkg.go.dev/github.com/pandatix/go-acm"><img src="https://shields.io/badge/-reference-blue?logo=go&style=for-the-badge" alt="reference"></a>
	<a href="https://goreportcard.com/report/github.com/pandatix/go-acm"><img src="https://goreportcard.com/badge/github.com/pandatix/go-acm?style=for-the-badge" alt="go report"></a>
	<a href="https://coveralls.io/github/pandatix/go-acm?branch=main"><img src="https://img.shields.io/coverallsCoverage/github/pandatix/go-acm?style=for-the-badge" alt="Coverage Status"></a>
	<a href=""><img src="https://img.shields.io/github/license/pandatix/go-acm?style=for-the-badge" alt="License"></a>
	<br>
	<a href="https://github.com/pandatix/go-acm/actions/workflows/ci.yaml"><img src="https://img.shields.io/github/actions/workflow/status/pandatix/go-acm/ci.yaml?style=for-the-badge&label=CI" alt="CI"></a>
	<a href="https://github.com/pandatix/go-acm/actions/workflows/codeql-analysis.yaml"><img src="https://img.shields.io/github/actions/workflow/status/pandatix/go-acm/codeql-analysis.yaml?style=for-the-badge&label=CodeQL" alt="CodeQL"></a>
	<br>
	<a href="https://securityscorecards.dev/viewer/?uri=github.com/pandatix/go-acm"><img src="https://img.shields.io/ossf-scorecard/github.com/pandatix/go-acm?label=openssf%20scorecard&style=for-the-badge" alt="OpenSSF Scoreboard"></a>
</div>

Go-ACM API scraps the ACM frontend to provide API functionality in Go. It currently covers:
- `action/doSearch`

## How to use

Examples use cases could be found in the [examples directory](examples).

The basic idea is to instanciate an `*ACMClient` and use it to issue calls that will be scrapped.
