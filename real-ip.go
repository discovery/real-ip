package real_ip

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type RemoteAddr struct {
	trusted_ipnets []net.IPNet
}

func cidrs_to_ipnets(cidrs []string) ([]net.IPNet, error) {
	ip_nets := []net.IPNet{}
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(strings.TrimSpace(cidr))
		if err != nil {
			return nil, err
		}
		ip_nets = append(ip_nets, *network)
	}
	return ip_nets, nil
}

func NewRemoteAddr(cidrs []string) (*RemoteAddr, error) {
	ipnets, err := cidrs_to_ipnets(cidrs)
	if err != nil {
		return nil, err
	}

	return &RemoteAddr{ipnets}, nil
}

func (ra *RemoteAddr) trusted(ip net.IP) bool {
	for _, ip_net := range ra.trusted_ipnets {
		if ip_net.Contains(ip) {
			return true
		}
	}
	return false
}

func (ra *RemoteAddr) GetIP(req *http.Request) (net.IP, error) {
	ips := strings.Split(req.Header.Get("X-Forwarded-For"), ",")
	for l := len(ips) - 1; l >= 0; l-- {
		ip := net.ParseIP(strings.TrimSpace(ips[l]))
		if ip != nil && !ra.trusted(ip) {
			return ip, nil
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		ip := net.ParseIP(host)
		if ip != nil {
			return ip, nil
		}
		err = errors.New(fmt.Sprintf(`Invalid remote address "%s" in request`, host))
	}

	return nil, err
}
