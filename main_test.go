package godnszone

import (
	"fmt"
	"testing"

	"github.com/rombintu/godnszone/utils"
)

func TestExecCommand(t *testing.T) {
	params := []string{"example.com", "example.com"}
	command := "named-checkzone"
	stdOut, err := utils.ExecCommand(command, params...)
	fmt.Println(stdOut)
	if err != nil {
		t.Error(err)
	}
}

func TestNewZoneChecker(t *testing.T) {
	zc := NewZoneChecker()
	fmt.Println(zc.IsValid("example.com", "example.com"))
	fmt.Println(zc.Output)
	fmt.Println(zc.Error)
	if zc.Error != nil {
		t.Error(zc.Error)
	}
}

func TestNewZoneReloader(t *testing.T) {
	zr := NewZoneReloader()
	fmt.Println(zr.Reload("example.com"))
	fmt.Println(zr.Output)
	fmt.Println(zr.Error)
	if zr.Error != nil {
		t.Error(zr.Error)
	}
}
