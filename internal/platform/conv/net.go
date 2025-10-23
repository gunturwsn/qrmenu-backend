package conv

import (
	"strconv"
	"strings"
)

func SplitHostPort(addr string, defaultHost string, defaultPort int) (string, int) {
	host := defaultHost
	port := defaultPort

	if addr == "" {
		return host, port
	}

	if strings.Count(addr, ":") == 0 {
		// hanya host
		return addr, port
	}

	h, p, found := strings.Cut(addr, ":")
	if !found {
		return host, port
	}
	host = h
	if pn, err := strconv.Atoi(p); err == nil {
		port = pn
	}
	return host, port
}
