# Exponent - Exponential Backoff for Go

[![Release](https://img.shields.io/github/release/tigraine/exponent.svg?style=flat-square)](https://github.com/tigraine/exponent/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![GoDoc](https://godoc.org/github.com/tigraine/exponent?status.svg)](http://godoc.org/github.com/tigraine/exponent)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigraine/exponent)](https://goreportcard.com/report/github.com/tigraine/exponent)

Library that provides various exponential backoff functions for your code.

Implementation of various exponential backoff strategies as outlined and tested in the AWS Blog: https://aws.amazon.com/de/blogs/architecture/exponential-backoff-and-jitter/
Supported strategies:
- Linear
- Exponential
- Full Jitter (exponential + random jitter)
- Decorrelated Jitter (exponential + partial random jitter)

## Usage

``` go
import (
  "os"

  "github.com/tigraine/exponent"
)

func main() {
  ctx := context.TODO()
  e := exponent.NewBackoff(12).WithContext(ctx)
  for e.Do() {
    // Do your work
    err := Work()
    if err != nil {
      e.Wait() // Sleep for the backoff
      continue
    }
    e.Success()
    break
  }
  if e.Failed() {
    os.Exit(1)
  }
}
```
