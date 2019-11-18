# Exponent - Exponential Backoff for Go

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
  e := exponent.NewBackoff(12)
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
