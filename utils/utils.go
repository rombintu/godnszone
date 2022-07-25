package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/miekg/dns"
)

// Execute any commands from shell
func ExecCommand(command string, params ...string) (string, error) {
	cmd := exec.Command(command, params...)
	var output string
	out, err := cmd.Output()
	if err != nil {
		return output, err
	}

	return string(out), nil
}

// Precheck file path
func ToValidPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(workDir, path)
	}
}

func ToValidName(name string) string {
	slice := strings.Split(name, "")
	if slice[len(slice)-1] == "." {
		return name
	} else {
		return name + "."
	}
}

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err := io.Copy(destination, source); err != nil {
		return err
	}
	return nil
}

func FilePathToBackupPath(fileName string) string {
	return fmt.Sprintf("%s.bak", fileName)
}

func ToFile(content, filePath string) error {
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func ToIsDomain(domain, name string) string {
	if domain == name {
		return "@"
	} else {
		return domain
	}
}

func ToIsTTL(globalTTL, ttl uint32) string {
	if ttl != globalTTL {
		return fmt.Sprint(ttl)
	} else {
		return ""
	}
}

func ToHeader(h string, s interface{}) string {
	return fmt.Sprintf("$%s %s\n", h, s)
}

func ToSubHeader(h string, s interface{}) string {
	return fmt.Sprintf("@%s %s\n", h, s)
}

func ToSOA(soa dns.SOA, newSerial uint32) string {
	return fmt.Sprintf(
		"@	%s	%s %s %s (\n\t%d ; serial\n\t%d ; retry\n\t%d ; expire\n\t%d ; minimum\n)\n",
		dns.ClassToString[soa.Hdr.Class],
		dns.TypeToString[soa.Hdr.Rrtype],
		soa.Ns,
		soa.Mbox,
		newSerial,
		soa.Retry,
		soa.Expire,
		soa.Minttl,
	)
}
