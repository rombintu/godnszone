package utils

import (
	"fmt"
	"runtime"
)

var (
	reset           = "\033[0m"
	red             = "\033[31m"
	green           = "\033[32m"
	ColorErr string = red + "ERROR" + reset
	ColorSuc string = green + "SUCCESS" + reset
)

const (
	ZoneIsValid    string = "File is valid"
	ZoneIsNotValid string = "File is not valid"

	SerialUpdated    string = "Serial updated"
	SerialNotUpdated string = "Srerial not updated"
	SerialFormat     string = "%04d%02d%02d00"

	RecordCreate   string = "Record created:"
	RecordUpdate   string = "Record updated:"
	RecordDelete   string = "Record deleted:"
	RecordNotFound string = "Record not found"
	RecordIsExists string = "Record is exists"
)

// Windows not supported color outputs
func ToOutput(output ...string) string {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		green = ""
	}
	return fmt.Sprint(output)
}
