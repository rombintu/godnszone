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
	filePath := utils.ToValidPath(fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	zw := newZoneWorker(filePath)
	zp := dns.NewZoneParser(bytes.NewReader(content), zoneName, zw.FilePath)

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		switch rr.Header().Rrtype {
		case dns.TypeSOA:
			zw.Zone.SOA = rr.(*dns.SOA)
			zw.Zone.Domain = zw.Zone.SOA.Header().Name
			zw.Zone.Serial = zw.Zone.SOA.Serial
			zw.Zone.Origin = zw.Zone.SOA.Hdr.Name
		default:
			rType := dns.TypeToString[rr.Header().Rrtype]
			zw.Zone.Records[rType] = append(zw.Zone.Records[rType], newExRRFromRR(rr, zp.Comment()))
		}
	}

	if err := zp.Err(); err != nil {
		zw.addError(err.Error())
	}

	return zw
}

func newExRRFromRR(rr dns.RR, comment string) ExRR {
	return ExRR{
		RR:      rr,
		Comment: comment,
	}
}

func newExRRFromString(rName, rType, rIP, comment string) (ExRR, error) {
	RR, err := dns.NewRR(fmt.Sprintf("%s %s %s", rName, rType, rIP))
	if err != nil {
		return ExRR{}, err
	}

	return ExRR{
		RR:      RR,
		Comment: comment,
	}, nil
}

func newExRRFromDry(rr, comment string) (ExRR, error) {
	RR, err := dns.NewRR(rr)
	if err != nil {
		return ExRR{}, err
	}

	return ExRR{
		RR:      RR,
		Comment: comment,
	}, nil
}

func TypeFromRR(rr ExRR) string {
	return dns.TypeToString[rr.RR.Header().Rrtype]
}

func ToRR(rr ExRR, domain string, gTTL uint16) string {
	origin := utils.ToIsDomain(domain, rr.RR.Header().Name)
	ttl := if 
	return fmt.Sprintf(
		"%s\t%s\t%s\t%s",
		origin,
		dns.ClassToString[rr.RR.Header().Class],
		dns.TypeToString[rr.RR.Header().Rrtype],

	)
}
