package godnszone

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
	"github.com/rombintu/godnszone/utils"
)

const (
	zoneTestName string = "example.com"
	fileTestName string = "example.com"
)

type testPairBool struct {
	values  []string
	average bool
}

func TestExecCommand(t *testing.T) {
	params := []string{zoneTestName, fileTestName}
	command := "named-checkzone"
	stdOut, err := utils.ExecCommand(command, params...)
	fmt.Println(stdOut)
	if err != nil {
		t.Error(err)
	}
}

func TestNewZoneChecker(t *testing.T) {
	zc := NewZoneChecker()
	fmt.Println(zc.IsValid(zoneTestName, fileTestName))
	fmt.Println(zc.Output)
	fmt.Println(zc.Error)
	if zc.Error != nil {
		t.Error(zc.Error)
	}
}

func TestNewZoneReloader(t *testing.T) {
	zr := NewZoneReloader()
	fmt.Println(zr.Reload(zoneTestName))
	fmt.Println(zr.Output)
	fmt.Println(zr.Error)
	if zr.Error != nil {
		t.Error(zr.Error)
	}
}

func TestNewZoneWorker(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	// fmt.Printf("%+v \n", zw.Zone.Records["A"])
	for i, r := range zw.Zone.Records {
		fmt.Println(i, r)
	}
	if zw.Errors != nil {
		t.Error(zw.Errors)
	}

}

func TestNewSerial(t *testing.T) {
	newSerial, err := newSerial(2022072101)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(newSerial)
}

func TestDeleteRecord(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	zw.delRecordByName("ns6.example.com", "TXT")
	fmt.Println(zw.getActions())
}

func TestAddRecord(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	rr, _ := dns.NewRR("ns6.example A 192.199.22.1")
	zw.addRecord(newExRR(rr, "Example create"))
	fmt.Println(zw.getActions())
}

func TestVerifyExist(t *testing.T) {
	rr1, _ := addDryRR("example.com NS ns1.example.com", "")
	rr2, _ := addDryRR("ns5.example.com A 192.26.238.38", "")
	rr3, _ := addDryRR("ns1.example.ru TXT fsdffsdf", "")
	var tests = []struct {
		in1  ExRR
		want bool
	}{
		{rr1, true},
		{rr2, true},
		{rr3, false},
	}
	zw := ZoneFromFile(zoneTestName, fileTestName)

	for _, tt := range tests {
		testName := fmt.Sprintf("%s", tt.in1)
		t.Run(testName, func(t *testing.T) {
			ans := zw.VerifyExist(tt.in1)
			if ans != tt.want {
				t.Errorf("Got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestVerifyExistByName(t *testing.T) {
	var tests = []struct {
		nameIn string
		typeIn string
		want   bool
	}{
		{"example.com", "NS", true},
		{"ns5.example.com", "A", true},
		{"ns1.example.ru", "TXT", false},
	}
	zw := ZoneFromFile(zoneTestName, fileTestName)

	for _, tt := range tests {
		testName := fmt.Sprintf("%s %s", tt.nameIn, tt.typeIn)
		t.Run(testName, func(t *testing.T) {
			ans := zw.VerifyExistByName(tt.nameIn, tt.typeIn)
			if ans != tt.want {
				t.Errorf("Got %t, want %t", ans, tt.want)
			}
		})
	}
}
