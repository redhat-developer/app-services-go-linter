# go-i18n-linter

[![Go Report Card](https://goreportcard.com/badge/github.com/alexal/go-i18n-linter)](https://goreportcard.com/report/github.com/alexal/go-i18n-linter)

**go-i18n-linter** plugin analyze source tree of Go files and validates the availability of i18n strings in *.toml files.
As of right now the project created exclusively for [Command Line Interface for RHOAS](https://github.com/redhat-developer/app-services-cli), however, 
you can adopt it for your needs by specifying the following command line options:
```bash
-path string
  Path to the directory with localization files. If nothing specified, linter will try to load i18n messages from files located in pkg/localize/locales directory.

-mustLocalize string
  Name of the function that loads an i18n message. (default "MustLocalize")

-mustLocalizeError string
  Name of the function that creates new error with i18n message. (default "MustLocalizeError")
```