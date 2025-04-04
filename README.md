# go-semver

A lightweight and idiomatic Go package to parse and compare semantic versions based on the [Semantic Versioning 2.0.0](https://semver.org/) specification.

## ✨ Features

- Parse semantic version strings (`MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`)
- Compare versions with full support for precedence rules
- Handle pre-release and build metadata
- Determine ordering between two versions
- Idiomatic Go API

## 📦 Installation

```bash
go get github.com/Nidal-Bakir/go-semver
```

## 🛠️ Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/go-semver"
)

func main() {
    v1, _ := semver.Parse("1.2.3-alpha.1")
    v2, _ := semver.Parse("1.2.3")

    if v1.IsLess(v2) {
        fmt.Println(v1, "is less than", v2)
    }
}
```

## 🔍 Version Comparison Rules

This package follows the precedence rules defined in the Semantic Versioning 2.0.0 spec, including:

* Precedence is determined by the major, minor, and patch versions, in that order.

* Pre-release versions have lower precedence than the associated normal version.

* Build metadata does not affect version precedence.

## ✅ Example Comparisons

| Version A        | Version B         | Result  |
|------------------|-------------------|---------|
| 1.0.0            | 2.0.0             | A < B   |
| 1.0.0-alpha      | 1.0.0             | A < B   |
| 1.0.0-alpha      | 1.0.0-alpha.1     | A < B   |
| 1.2.3+build1     | 1.2.3+build2      | A == B  |
