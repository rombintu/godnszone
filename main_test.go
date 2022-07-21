package godnszone

import (
	"fmt"
	"testing"

	"github.com/rombintu/godnszone/utils"
)

const (
	zoneTestName string = "example.com"
	fileTestName string = "example.com"
)

func TestExecCommand(t *testing.T) {
	params := []string{zoneTestName, fileTestName}
	command := "named-checkzone"
	stdOut, err := utils.ExecCommand(command, params...)
	fmt.Println(stdOut)
	if err != nil {
		t.Error(err)
	}
}

func TestNewZoneChecker(t *testing.T) {
	zc := NewZoneChecker()
	fmt.Println(zc.IsValid(zoneTestName, fileTestName))
	fmt.Println(zc.Output)
	fmt.Println(zc.Error)
	if zc.Error != nil {
		t.Error(zc.Error)
	}
}

func TestNewZoneReloader(t *testing.T) {
	zr := NewZoneReloader()
	fmt.Println(zr.Reload(zoneTestName))
	fmt.Println(zr.Output)
	fmt.Println(zr.Error)
	if zr.Error != nil {
		t.Error(zr.Error)
	}
}

func TestNewZoneWorker(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	// fmt.Printf("%+v \n", zw.Zone.Records["A"])
	for i, r := range zw.Zone.Records {
		fmt.Println(i, r)
	}
	if zw.Errors != nil {
		t.Error(zw.Errors)
	}

}

// func TestTable(t *testing.T) { // TODO
// 	zw := ZoneFromFile(zoneTestName, fileTestName)
// 	if zw.Errors != nil {
// 		t.Error(zw.Errors)
// 	}
// 	zw.Table()
// }

func TestNewSerial(t *testing.T) {
	newSerial, err := NewSerial(2022072101)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(newSerial)
}
