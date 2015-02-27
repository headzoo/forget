Forget
======
Forget is a Go client library used to interact with [Forgettable](https://github.com/bitly/forgettable) servers.

[![Build Status](https://img.shields.io/travis/headzoo/forget/master.svg)](https://travis-ci.org/headzoo/forget)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/headzoo/forget)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/headzoo/forget/master/LICENSE.md)

### Installation
Download the library using go.  
`go get github.com/headzoo/forget`

Import the library into your project.  
`import "github.com/headzoo/forget"`


### General Usage
First make sure your Forgettable server is up and running. In this example our server is running at ficticious URL
http://forgettable.io:51000.

```go
package main

import (
  "github.com/headzoo/forget"
  "fmt"
)

func main() {
  client := forget.NewClient("http://forgettable.io:51000")
  
  err := client.Increment("colors", "red")
  if err != nil {
    panic(err)
  }
  
  res, err := client.Distribution("colors")
  if err != nil {
    panic(err)
  }
  
  fmt.Printf("Got status code %d for distribution %s.\n", res.StatusCode, res.Distribution.Name)
  for _, value := range res.Distribution.Values {
    fmt.Printf("Field %s has a probability of %f.\n", value.Field, value.Probability)
  }
  
  res, err = client.MostProbable("colors", 10)
  if err != nil {
    panic(err)
  }
  
  fmt.Println("Displaying fields from most probable to least.")
  for i, value := res.Distribution.Values {
    fmt.Printf("%d. %s\n", i, value.Field)
  }
}
```
