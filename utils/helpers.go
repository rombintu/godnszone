package utils

import (
	"fmt"
	"runtime"
)

var (
	reset           = "\033[0m"
	red             = "\033[31m"
	green           = "\033[32m"
	ColorErr string = red + "FAILED" + reset
	ColorSuc string = green + "SUCCESS" + reset
)

const (
	ZoneIsValid    string = "file is valid"
	ZoneIsNotValid string = "file is not valid"

	SerialUpdated    string = "serial updated"
	SerialNotUpdated string = "srerial not updated"
	SerialFormat     string = "%04d%02d%02d00"

	RecordCreate    string = "record created:"
	RecordNotCreate string = "record not created:"

	RecordUpdate    string = "record updated:"
	RecordNotUpdate string = "record not updated:"

	RecordDelete    string = "record deleted:"
	RecordNotDelete string = "record not deleted:"

	RecordNotFound string = "record not found"
	RecordIsExists string = "record is exists"
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
