package godnszone

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
)

const (
	zoneName string = "example.com"
	fileName string = "example.com"
)

type ExRR struct {
	RR      dns.RR
	Comment string
}

type Zone struct {
	SOA     *dns.SOA
	Domain  string
	Serial  uint32
	Records map[string][]ExRR
}

type ZoneWorker struct {
	// Parser      *dns.ZoneParser
	Zone   *Zone
	Errors []error
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

func (zw *ZoneWorker) TableByType(t string) {
	// TODO
}

func ZoneFromFile(zoneName, fileName string) *ZoneWorker {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("READFILE: ", err)
	}

	zw := newZoneWorker()
	zp := dns.NewZoneParser(bytes.NewReader(file), zoneName, fileName)

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		switch rr.Header().Rrtype {
		case dns.TypeSOA:
			zw.Zone.SOA = rr.(*dns.SOA)
			zw.Zone.Domain = rr.(*dns.SOA).Header().Name
			zw.Zone.Serial = rr.(*dns.SOA).Serial
		default:

			recordType := dns.TypeToString[rr.Header().Rrtype]
			zw.Zone.Records[recordType] = append(zw.Zone.Records[recordType], newExRR(rr, zp.Comment()))
		}
	}

	if err := zp.Err(); err != nil {
		zw.Errors = append(zw.Errors, err)
	}

	return zw
}

func newExRR(rr dns.RR, comment string) ExRR {
	return ExRR{
		RR:      rr,
		Comment: comment,
	}
}

func newRR(record, typ, ip, comment string) (ExRR, error) {
	RR, err := dns.NewRR(fmt.Sprintf("%s %s %s", record, typ, ip))
	if err != nil {
		return ExRR{}, err
	}

	return ExRR{
		RR:      RR,
		Comment: comment,
	}, nil
}

func main() {
	rr, err := newRR("bfd", "A", "192.168.10.20", "запись для проекта")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rr)

}
