package godnszone

import (
	"github.com/miekg/dns"
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
}

func (zw *ZoneWorker) delRecordByName(rName, rType string) {

}

func (zw *ZoneWorker) Save() {

}

func (zw *ZoneWorker) UpdateSerial() error {
	newSerial, err := newSerial(zw.Zone.SOA.Serial)
	if err != nil {
		zw.Errors = append(zw.Errors, err)
		return err
	}
	zw.Zone.SOA.Serial = newSerial
	return nil
}
