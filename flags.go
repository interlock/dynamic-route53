package main

import (
	"flag"
	"fmt"
	"strings"
)

var validNetworkFlags = []string{"ipv4", "ipv6"}
var networkFlag = flag.String("networks", "ipv4", "ipv4 or ipv6, sets A or AAAA as appropriate")
var validLookupFlags = []string{"ifconfig.co"}
var lookupFlag = flag.String("lookup", "ifconfig.co", "ifconfig.co for lookup")
var domainFlag = flag.String("domain", "", "domain to update on Route53")
var hostedZoneIDFlag = flag.String("hosted-zone-id", "", "hosted zone id on Route53")
var ttlFlag = flag.Int64("ttl", 3600, "dns record ttl")
var commentFlag = flag.String("comment", "", "comment to set on the Route53 Record")

func containsValue(value string, options []string) bool {
	for _, i := range options {
		if strings.Compare(value, i) == 0 {
			return true
		}
	}
	return false
}

func validateFlags() []error {
	var errors []error

	if containsValue(*networkFlag, validNetworkFlags) == false {
		errors = append(errors, fmt.Errorf("Expected network flag to be %s", strings.Join(validNetworkFlags, ",")))
	}
	if containsValue(*lookupFlag, validLookupFlags) == false {
		errors = append(errors, fmt.Errorf("Expected lookup flag %s", strings.Join(validLookupFlags, ",")))
	}
	if len(*domainFlag) == 0 {
		errors = append(errors, fmt.Errorf("domain is required"))
	}
	if len(*hostedZoneIDFlag) == 0 {
		errors = append(errors, fmt.Errorf("hosted-zone-id is required"))
	}

	return errors
}
