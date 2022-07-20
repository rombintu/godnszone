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
	fmt.Printf("%+v \n", zw.Zone.Records)
	if zw.Errors != nil {
		t.Error(zw.Errors)
	}

}
