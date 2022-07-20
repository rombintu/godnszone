package godnszone

import (
	"github.com/rombintu/godnszone/utils"
)

type ZoneChecker struct {
	CheckZone string
	Output    string
	Error     error
}

func NewZoneChecker() *ZoneChecker {
	return &ZoneChecker{
		CheckZone: "named-checkzone",
	}
}

func (z *ZoneChecker) IsValid(zoneName, fileName string) bool {
	params := []string{
		zoneName,
		fileName,
	}

	stdOut, err := utils.ExecCommand(z.CheckZone, params...)
	if err != nil {
		z.Error = err
		return false
	}
	z.Output = stdOut
	return true
}
