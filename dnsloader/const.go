package dnsloader

// ReturnCode for return code of each call function
type ReturnCode int

//dns loader running status
const (
	STATUS_ORIGINAL uint32 = 0
	STATUS_STARTING uint32 = 0
	STATUS_STARTED  uint32 = 0
	STATUS_STOPPING uint32 = 0
	STATUS_STOPPED  uint32 = 0
)
const (
	CALL_NOT_FINISH   uint32 = 0
	CALL_SUCCESS_DONE uint32 = 1
	CALL_TIMEOUT      uint32 = 2
)

// Return Code
const (
	RetNoError   ReturnCode = 0
	RetServfail             = 0x2
	RetFormError            = 0x1
	RetNxDomain             = 0x3
	RetNoImp                = 0x4
	RetRefused              = 0x5
	RetTimeout              = 0xf0
)

const (
	// valid dnsRR_Header.Rrtype and dnsQuestion.qtype
	dnsTypeA     = 1
	dnsTypeNS    = 2
	dnsTypeMD    = 3
	dnsTypeMF    = 4
	dnsTypeCNAME = 5
	dnsTypeSOA   = 6
	dnsTypeMB    = 7
	dnsTypeMG    = 8
	dnsTypeMR    = 9
	dnsTypeNULL  = 10
	dnsTypeWKS   = 11
	dnsTypePTR   = 12
	dnsTypeHINFO = 13
	dnsTypeMINFO = 14
	dnsTypeMX    = 15
	dnsTypeTXT   = 16
	dnsTypeAAAA  = 28
	dnsTypeSRV   = 33
	dnsTypeAXFR  = 252
	dnsTypeMAILB = 253
	dnsTypeMAILA = 254
	dnsTypeALL   = 255

	// valid dnsQuestion.qclass
	dnsClassINET   = 1
	dnsClassCSNET  = 2
	dnsClassCHAOS  = 3
	dnsClassHESIOD = 4
	dnsClassANY    = 255

	// dnsMsg.rcode
	dnsRcodeSuccess        = 0
	dnsRcodeFormatError    = 1
	dnsRcodeServerFailure  = 2
	dnsRcodeNameError      = 3
	dnsRcodeNotImplemented = 4
	dnsRcodeRefused        = 5
)

// DNSType define the dns query type for packet build
var DNSType = map[string]uint16{
	"A":     1,
	"NS":    2,
	"MD":    3,
	"MF":    4,
	"CNAME": 5,
	"SOA":   6,
	"MB":    7,
	"MG":    8,
	"MR":    9,
	"NULL":  10,
	"WKS":   11,
	"PTR":   12,
	"HINFO": 13,
	"MINFO": 14,
	"MX":    15,
	"TXT":   16,
	"AAAA":  28,
	"SRV":   33,
	"AXFR":  252,
	"MAILB": 253,
	"MAILA": 254,
	"ALL":   255,
}

// QClass define the dns query class
var QClass = map[string]int{
	"INET":   1,
	"CSNET":  2,
	"CHAOS":  3,
	"HESIOD": 4,
	"ANY":    255,
}

// DNSRcode define the dns return code
var DNSRcode = map[string]int{
	"Success":        0,
	"FormatError":    1,
	"ServerFailure":  2,
	"NameError":      3,
	"NotImplemented": 4,
	"Refused":        5,
	"YXDOMAIN":       6,
	"XRRSET":         7,
	"NotAuth":        8,
	"NotInZone":      9,
}

// DNSRcodeReverse define the real code to string map
var DNSRcodeReverse = map[uint8]string{
	0: "Success",
	1: "FormatError",
	2: "ServerFail",
	3: "NXDOMAIN",
	4: "NotImplemented",
	5: "Refused",
	6: "YXDOMAIN",
	7: "XRRSET",
	8: "NotAuth",
	9: "NotInZone",
}
