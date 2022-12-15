package api

import (
	"fmt"
	"net"
	"net/http"
)

// check what will be as remote addr
// x-forwarded-for will keep the addresses in case of proxies etc etc.
func extractIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("wrong remote addr: %s: %w", r.RemoteAddr, err)
	}

	return ip, nil
}
