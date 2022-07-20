# DNS Zone tools

### Description 
Alternative Golang module for the [dnszone](https://pypi.org/project/dnszone/) (python-package)

### Dependencies
For DnsZone
* TODO

For Check/Reload
* bind
* bind9tools

## Install
```bash
$ go get "github.com/rombintu/godnszone"
```

## Usage Check/Reload zone
```go
package main

import (
	"fmt"

	dnszone "github.com/rombintu/godnszone"
)

const (
	zoneName string = "example.com"
	fileName string = "example.com"
)

// Example
func main() {
	zc := dnszone.NewZoneChecker()
	zr := dnszone.NewZoneReloader()

	// If you get errors, try change paths
	zc.CheckZone = "/usr/bin/named-checkzone" // Default: "named-checkzone"
	zr.ReloadZone = "/usr/sbin/rndc"          // Default: "rndc"

	if zc.IsValid(zoneName, fileName) {
		fmt.Println("Checker:", zc.Output)
		if zr.Reload(zoneName) {
			fmt.Println("Reloader", zr.Output)
		} else {
			fmt.Println("Reloader:", zr.Error)
		}
	} else {
		fmt.Println("Checker: ", zc.Error)
	}
}
```