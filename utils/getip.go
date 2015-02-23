package utils

import (
	"net"
)

//If this returns an error, just return the input string
func IpFromRemoteAddr(remote string) string {
	out, _, err := net.SplitHostPort(remote)
	if err != nil {
		return remote
	}
	return out
}
