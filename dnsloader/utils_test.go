package dnsloader

import (
	"testing"
)

func TestGetDNSTypeCodeFromString(t *testing.T) {
	var testCase = []struct {
		input  string
		expect int
	}{{
		input:  "a",
		expect: 1,
	}, {
		input:  "A",
		expect: 1,
	}, {
		input:  "AAAA",
		expect: 28,
	}, {
		input:  "aaaa",
		expect: 28,
	}, {
		input:  "not-exist",
		expect: 0,
	},
	}
	for _, obj := range testCase {
		output := GetDNSTypeCodeFromString(obj.input)
		if output != obj.expect {
			t.Errorf("Expect %d: Got %d", obj.expect, output)
		}
	}
}

func TestGetDetailInfo(t *testing.T) {
	var detail string
	var testCase = []struct {
		code   ReturnCode
		detail string
	}{{
		code:   RET_FORMERR,
		detail: "Query Format Error",
	}, {
		code:   RET_NO_ERROR,
		detail: "Success",
	}, {
		code:   RET_NOTIMP,
		detail: "Query Type Not Support",
	}, {
		code:   RET_NXDOMAIN,
		detail: "Query Not Exist",
	}, {
		code:   RET_REFUSED,
		detail: "Query Refused",
	}, {
		code:   RET_SERVFAIL,
		detail: "Server Fail",
	}, {
		code:   100001,
		detail: "Unknown",
	}}
	for _, obj := range testCase {
		detail = GetDetailInfo(obj.code)
		if detail != obj.detail {
			t.Errorf("Expect %s: Got %s", obj.detail, detail)
		}
	}
}
