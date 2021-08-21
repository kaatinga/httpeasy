[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/httpeasy/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/httpeasy/branch/main/graph/badge.svg)](https://codecov.io/gh/kaatinga/httpeasy)
[![lint workflow](https://github.com/kaatinga/httpeasy/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/httpeasy/actions?query=workflow%3Alinter)
[![Help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/httpeasy/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# httpeasy
The package provides everything to create a simple web service. You merely have to announce handlers.

## 1. Installation

Use go get.

	go get github.com/kaatinga/httpeasy

Then import the validator package into your own code.

	import "github.com/kaatinga/httpeasy"

## 2. Usage

Prepare a function that complies with `SetUpHandlers` type. The function should contain some routes, for example.

    func SetUpHandlers(r *httprouter.Router, _ *sql.DB) {
	    r.GET("/", Welcome)
    }

The package contains a ready config model, set field values in that structure, easies way is to use [settings](https://github.com/kaatinga/settings) package:

    err := settings.LoadSettings(&config.Config)
    if err != nil {
        ...
    }

Run your server:

    err = config.Config.Launch(SetUpHandlers, logger)
    if err != nil {
        ...
    }