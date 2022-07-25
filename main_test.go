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
	fmt.Println("TEST 1 >>")
	fmt.Println(zw.Zone.Records["TXT"])
	if err := zw.DeleteRecordByName("_acme-challenge.brr.example.com", "TXT"); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())
	fmt.Println(zw.Zone.Records["TXT"])

	fmt.Println("TEST 2 >>")
	fmt.Println(zw.Zone.Records["A"])
	if err := zw.DeleteRecordByName("site.example.com", "A"); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())
	fmt.Println(zw.Zone.Records["A"])
}

func TestAddRecord(t *testing.T) {
	fmt.Println("TEST WITH NOT EXIST >>")
	zw := ZoneFromFile(zoneTestName, fileTestName)
	rr1, _ := dns.NewRR("ns6.example A 192.199.228.1")
	if err := zw.AddRecord(NewExRRFromRR(rr1, "Example create")); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())

	fmt.Println("TEST WITH EXIST >>")
	rr2, _ := dns.NewRR("ns5.example.com A 192.26.238.38")
	if err := zw.AddRecord(NewExRRFromRR(rr2, "Example create")); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())
}

func TestUpdateRecord(t *testing.T) {
	fmt.Println("TEST WITH NOT EXIST >>")
	zw := ZoneFromFile(zoneTestName, fileTestName)
	rr1, _ := dns.NewRR("ns122.example A 192.199.228.1")
	if err := zw.UpdateRecordByName(
		"ns50.example.com",
		"ABV",
		NewExRRFromRR(rr1, ""),
	); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())

	fmt.Println("TEST WITH EXIST >>")
	rr2, _ := dns.NewRR("ns5.example.com A 192.26.238.38")
	if err := zw.UpdateRecordByName(
		"ns5.example.com",
		"A",
		NewExRRFromRR(rr2, "Example create"),
	); err != nil {
		t.Error(err)
	}
	fmt.Println(zw.GetActions())
}

func TestVerifyExist(t *testing.T) {
	rr1, _ := NewExRRFromDry("example.com NS ns1.example.com", "")
	rr2, _ := NewExRRFromDry("ns5.example.com A 192.26.238.38", "")
	rr3, _ := NewExRRFromDry("ns1.example.ru TXT fsdffsdf", "")
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

func TestToValidName(t *testing.T) {
	var tests = []struct {
		in, want string
	}{
		{"some", "some."},
		{"some.", "some."},
		{"string..", "string."},
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("%s %s", tt.in, tt.want)
		t.Run(testName, func(t *testing.T) {
			ans := utils.ToValidName(tt.in)
			if ans != tt.want {
				t.Errorf("Got: %s\n Want: %s", ans, tt.want)
			}
		})
	}
}

func TestBackup(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	if err := zw.Backup(); err != nil {
		t.Error(err)
	}
}

func TestGetSOA(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	fmt.Printf("%+v", *zw.Zone.SOA)
}

func TestSave(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	zw.Save(true)
}

func TestGetOneRRA(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	fmt.Println(utils.GetPayloadFromRR(zw.Zone.Records["A"][0].RR))
}

func TestGetOneRRTXT(t *testing.T) {
	zw := ZoneFromFile(zoneTestName, fileTestName)
	fmt.Println(utils.GetPayloadFromRR(zw.Zone.Records["TXT"][0].RR))
}
