package dnsloader

import (
	"math/rand"
	"strings"
	"time"
)

var letters = []byte("1234567890abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetDetailInfo
func GetDetailInfo(code ReturnCode) string {
	var detail string
	switch code {
	case RetNoError:
		detail = "Success"
	case RetFormError:
		detail = "Query Format Error"
	case RetNoImp:
		detail = "Query Type Not Support"
	case RetNxDomain:
		detail = "Query Not Exist"
	case RetRefused:
		detail = "Query Refused"
	case RetServfail:
		detail = "Server Fail"
	default:
		detail = "Unknown"
	}
	return detail
}

// GetDNSTypeCodeFromString return the true code of dns type
func GetDNSTypeCodeFromString(typeString string) uint16 {
	queryString := strings.ToUpper(typeString)
	code, ok := DNSType[queryString]
	if ok {
		return code
	}
	return 0
}

//GenRandomDomain will generate the random domain name with the fix length
func GenRandomDomain(length int, domain string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + "." + domain
}
