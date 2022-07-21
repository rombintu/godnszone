package godnszone

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

// func (zw *ZoneWorker) Table() { // TODO

// 	w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', tabwriter.AlignRight)
// 	fmt.Fprintln(w, "TYPE\tRECORDS\tCOMMENTS\t")

// 	for _, rr := range zw.Zone.Records {
// 		buffLine := ""
// 		for _, r := range rr {
// 			buffLine += fmt.Sprintf("%s [%s]\t", r.RR.Header().Name, r.Comment)
// 		}
// 		fmt.Fprintln(w, buffLine)
// 	}
// 	w.Flush()
// }

func NewSerial(oldSerial uint32) (uint32, error) {
	t := time.Now().Local()
	parseUint, err := strconv.ParseUint(
		fmt.Sprintf("%04d%02d%02d00", t.Year(), t.Month(), t.Day()),
		10,
		32,
	)
	if err != nil {
		return 0, err
	}
	newSerial := uint32(parseUint)
	if newSerial <= oldSerial {
		newSerial = oldSerial + 1
	}
	return newSerial, nil
}

func (zw *ZoneWorker) Save() {

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
