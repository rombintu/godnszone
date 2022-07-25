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

## Usage Add/Update/Delete records
```go
package main

import (
	"fmt"
	"log"

	dnszone "github.com/rombintu/godnszone"
)

const (
	zoneName string = "example.com"
	fileName string = "example.com"
)

func main() {
	// Create ZoneWorker
	zw := dnszone.ZoneFromFile(zoneName, fileName)
	rr1, _ := dnszone.NewExRRFromDry("ns6.example A 192.199.228.1", "Some comment")
	if err := zw.AddRecord(rr1); err != nil {
		log.Println(err)
	}

	// Create the same record [is FAILED]
	rr2, _ := dnszone.NewExRRFromDry("ns6.example A 192.199.119.1", "Some comment 2")
	if err := zw.AddRecord(rr2); err != nil {
		log.Println(err)
	}

	// Update record
	if err := zw.UpdateRecordByName("ns6.example", "A", rr2); err != nil {
		log.Println(err)
	}

	// Delete record
	if err := zw.DeleteRecordByName("ns6.example", "A"); err != nil {
		log.Println(err)
	}

	// Print my actions
	actions := zw.GetActions()
	for i, a := range actions {
		fmt.Printf("%d) %s\n", i+1, a)
	}
}

```
## Output
```log
1) [record created: ns6.example.        3600    IN      A       192.199.228.1 SUCCESS]
2) [record is exists ns6.example.       3600    IN      A       192.199.119.1 FAILED]
3) [record updated: ns6.example A SUCCESS]
4) [record deleted: ns6.example A SUCCESS]
```