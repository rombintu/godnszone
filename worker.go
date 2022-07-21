package godnszone

import "github.com/miekg/dns"

type Zone struct {
	SOA     *dns.SOA
	Domain  string
	Serial  uint32
	Records map[string][]ExRR
}

// Main object for manage zones
type ZoneWorker struct {
	Zone     *Zone
	Precheck []string
	Errors   []error
}

func newZone() *Zone {
	return &Zone{
		Records: make(map[string][]ExRR),
	}
}

func newZoneWorker() *ZoneWorker {
	return &ZoneWorker{
		Zone: newZone(),
	}
}

func (zw *ZoneWorker) Save() {

}

func (zw *ZoneWorker) NewSerial() error {
	newSerial, err := newSerial(zw.Zone.SOA.Serial)
	if err != nil {
		zw.Errors = append(zw.Errors, err)
		return err
	}
	zw.Zone.SOA.Serial = newSerial
	return nil
}
