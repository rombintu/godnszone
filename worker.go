package godnszone

import (
	"errors"

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

func (zw *ZoneWorker) addError(err string) {
	zw.Errors = append(zw.Errors, errors.New(err))
}

func (zw *ZoneWorker) addAction(action string) {
	zw.Actions = append(zw.Actions, action)
}

func (zw *ZoneWorker) GetActions() []string {
	return zw.Actions
}

func (zw *ZoneWorker) AddRecord(rr ExRR) error {
	if zw.VerifyExist(rr) {
		zw.addAction(utils.ToOutput(utils.RecordIsExists, rr.RR.String(), utils.ColorErr))
		return errors.New(utils.ToOutput(utils.RecordNotCreate, rr.RR.String(), utils.ColorErr))
	}
	rType := dns.TypeToString[rr.RR.Header().Rrtype]
	zw.Zone.Records[rType] = append(zw.Zone.Records[rType], rr)
	zw.addAction(utils.ToOutput(utils.RecordCreate, rr.RR.String(), utils.ColorSuc))
	return nil
}

func (zw *ZoneWorker) DeleteRecordByName(rName, rType string) error {

	for i, rr := range zw.Zone.Records[rType] {
		if rr.RR.Header().Name == dns.CanonicalName(rName) {
			zw.Zone.Records[rType] = append(
				zw.Zone.Records[rType][:i],
				zw.Zone.Records[rType][i+1:]...,
			)
			zw.addAction(utils.ToOutput(utils.RecordDelete, rName, rType, utils.ColorSuc))
			return nil
		}
	}

	zw.addAction(utils.ToOutput(utils.RecordNotDelete, rName, rType, utils.ColorErr))
	return errors.New(utils.ToOutput(utils.RecordNotDelete, rName, rType, utils.ColorErr))
}

func (zw *ZoneWorker) UpdateRecordByName(rName, rType string, newRR ExRR) error {
	if !zw.VerifyExistByName(rName, rType) {
		zw.addAction(utils.ToOutput(utils.RecordNotUpdate, rName, rType, utils.ColorErr))
		return errors.New(utils.ToOutput(utils.RecordNotFound, rName, rType, utils.ColorErr))
	}
	for i, rr := range zw.Zone.Records[rType] {
		if rr.RR.Header().Name == dns.CanonicalName(rName) {
			zw.Zone.Records[rType][i] = rr
			zw.addAction(utils.ToOutput(utils.RecordUpdate, rName, rType, utils.ColorSuc))
			return nil
		}
	}
	return errors.New(utils.ToOutput(utils.RecordNotUpdate, rName, rType, utils.ColorErr))
}

// If record exist return TRUE
func (zw *ZoneWorker) VerifyExistByName(rName, rType string) bool {
	exist := false
	for _, rr := range zw.Zone.Records[rType] {
		if rr.RR.Header().Name == dns.CanonicalName(rName) {
			exist = true
		}
	}
	return exist
}

// If record exist return TRUE
func (zw *ZoneWorker) VerifyExist(rr ExRR) bool {
	exist := false
	for _, r := range zw.Zone.Records[TypeFromRR(rr)] {
		if r.RR.Header().Name == rr.RR.Header().Name {
			exist = true
		}
	}
	return exist
}

func (zw *ZoneWorker) Save(autoSerial bool) error {
	// TODO
	if autoSerial {
		if err := zw.UpdateSerial(); err != nil {
			return err
		}
	}
	// TODO
	return nil
}

// Create new file.bak (from zw.FilePath)
func (zw *ZoneWorker) Backup() error {
	if err := utils.Copy(
		zw.FilePath,
		utils.FilePathToBackupPath(zw.FilePath),
	); err != nil {
		return err
	}
	return nil
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
