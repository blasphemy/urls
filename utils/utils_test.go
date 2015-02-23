package utils

import (
	"testing"
)

func TestGetIp1(t *testing.T) {
	k := IpFromRemoteAddr("127.0.0.1:4339")
	if k != "127.0.0.1" {
		t.Fail()
	}
}

func TestGetIp2(t *testing.T) {
	k := IpFromRemoteAddr("Test")
	if k != "Test" {
		t.Fail()
	}
}
