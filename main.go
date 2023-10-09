package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type DomainEmailRecords struct {
	domain      string
	hasMx       bool
	hasSPF      bool
	spfRecord   string
	hasDMARC    bool
	dmarcRecord string
}

func main() {
	fmt.Sprintf("Hello World")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Here's what you'll get: ")
	fmt.Printf("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcRecord\n\n")
	fmt.Printf("Enter a domain: ")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		//log.Println(err)
		log.Fatal("Error: Could not read from input: %v\n", err)
	}

}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: Could not lookup MX records for %s: %v\n", domain, err)
	}

	if len(mxRecords) > 0 {
		hasMx = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: Could not lookup TXT records for %s: %v\n", domain, err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: Could not lookup DMARC record for %s: %v\n", domain, err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord)

}

func isValidDomain(domain string) bool {

	if !(strings.Contains(domain, ".")) {
		return false
	}

	if len(domain) > 255 {
		return false
	}

	if domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1]
	}
	for _, v := range strings.Split(domain, ".") {
		if len(v) > 63 {
			return false
		}
	}
	return true

}
