package godnszone

import (
	"github.com/rombintu/godnszone/utils"
)

type ZoneReloader struct {
	ReloadZone string
	Output     string
	Error      error
}

func NewZoneReloader() *ZoneReloader {
	return &ZoneReloader{
		ReloadZone: "rndc",
	}
}

func (z *ZoneReloader) Reload(zoneName string) bool {
	params := []string{
		zoneName,
	}

	stdOut, err := utils.ExecCommand(z.ReloadZone, params...)
	if err != nil {
		z.Error = err
		return false
	}
	z.Output = stdOut
	return true
}
