package godnszone

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/miekg/dns"
	"github.com/rombintu/godnszone/utils"
)

// Custom record with comment
type ExRR struct {
	RR      dns.RR
	Comment string
}

func newSerial(oldSerial uint32) (uint32, error) {
	t := time.Now().Local()
	parseUint, err := strconv.ParseUint(
		fmt.Sprintf(utils.SerialFormat, t.Year(), t.Month(), t.Day()),
		10,
		32,
	)
	if err != nil {
		return oldSerial, err
	}
	newSerial := uint32(parseUint)
	if newSerial <= oldSerial {
		newSerial = oldSerial + 1
	}
	return newSerial, nil
}

func ZoneFromFile(zoneName, fileName string) *ZoneWorker {
	file, err := os.ReadFile(utils.ToValidPath(fileName))
	if err != nil {
		log.Fatalf("%+v", err)
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

func addRR(name, t, ip, comment string) (ExRR, error) {
	RR, err := dns.NewRR(fmt.Sprintf("%s %s %s", name, t, ip))
	if err != nil {
		return ExRR{}, err
	}

	return ExRR{
		RR:      RR,
		Comment: comment,
	}, nil
}
