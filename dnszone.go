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
	dns.RR
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
	Zone   Zone
	Errors []error
}

func newZoneWorker() *ZoneWorker {
	return &ZoneWorker{}
}

func ZoneFromFile(zoneName, fileName string) *ZoneWorker {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("READFILE: ", err)
	}

	zw := newZoneWorker()
	zp := dns.NewZoneParser(bytes.NewReader(file), zoneName, fileName)

	for r, ok := zp.Next(); ok; r, ok = zp.Next() {
		switch r.Header().Rrtype {
		case dns.TypeSOA:
			zw.Zone.SOA = r.(*dns.SOA)
		default:
			// zw.Zone.Records[r.Header().Name]
			// newExRR(r, zp.Comment())
		}
	}

	if err := zp.Err(); err != nil {
		log.Println(err)
	}
	zw.Zone.Domain = zw.Zone.SOA.Header().Name
	zw.Zone.Serial = zw.Zone.SOA.Serial
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
