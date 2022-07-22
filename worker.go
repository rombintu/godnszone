package godnszone

import (
	"github.com/miekg/dns"
	"github.com/rombintu/godnszone/utils"
)

type Zone struct {
	SOA     *dns.SOA
	Domain  string
	Serial  uint32
	Records map[string][]ExRR
}

// Main object for manage zones
type ZoneWorker struct {
	Zone     *Zone
	FilePath string
	Actions  []string
	Errors   []error
}

func newZone() *Zone {
	return &Zone{
		Records: make(map[string][]ExRR),
	}
}

func newZoneWorker(filePath string) *ZoneWorker {
	return &ZoneWorker{
		Zone:     newZone(),
		FilePath: filePath,
	}
}

func (zw *ZoneWorker) addAction(action string) {
	zw.Actions = append(zw.Actions, action)
}

func (zw *ZoneWorker) getActions() []string {
	return zw.Actions
}

func (zw *ZoneWorker) addRecord(rr ExRR) {
	rType := dns.TypeToString[rr.RR.Header().Rrtype]
	zw.Zone.Records[rType] = append(zw.Zone.Records[rType], rr)
	zw.addAction(utils.ToOutput(utils.RecordCreate, rr.RR.String(), utils.ColorSuc))
}

func (zw *ZoneWorker) delRecordByName(rName, rType string) {
	deleted := false
	for i, rr := range zw.Zone.Records[rType] {
		if rr.RR.Header().Name == rName+"." {
			zw.Zone.Records[rType] = append(
				zw.Zone.Records[rType][:i],
				zw.Zone.Records[rType][i+1:]...,
			)
			deleted = true
			zw.addAction(utils.ToOutput(utils.RecordDelete, rName, rType, utils.ColorSuc))
		}
	}
	if !deleted {
		zw.addAction(utils.ToOutput(utils.RecordNotFound, rName, rType, utils.ColorErr))
	}
}

func (zw *ZoneWorker) VerifyExistByName(rName, rType string) bool {
	exist := false
	for _, rr := range zw.Zone.Records[rType] {
		if rr.RR.Header().Name == rName+"." {
			exist = true
		}
	}
	return exist
}

func (zw *ZoneWorker) VerifyExist(rr ExRR) bool {
	exist := false
	for _, r := range zw.Zone.Records[TypeFromRR(rr)] {
		if r.RR.Header().Name == rr.RR.Header().Name {
			exist = true
		}
	}
	return exist
}

func (zw *ZoneWorker) Save() {
	// TODO
}

func (zw *ZoneWorker) UpdateSerial() error {
	newSerial, err := newSerial(zw.Zone.SOA.Serial)
	if err != nil {
		zw.Errors = append(zw.Errors, err)
		zw.addAction(utils.ToOutput(utils.SerialNotUpdated, utils.ColorErr))
		return err
	}
	zw.Zone.SOA.Serial = newSerial
	zw.addAction(utils.ToOutput(utils.SerialUpdated, utils.ColorErr))
	return nil
}
