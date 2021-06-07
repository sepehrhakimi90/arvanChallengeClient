package utils

import (
	"net"
)

func GetDomainIPv4 (host string) ([]net.IP, error){
	result := make([]net.IP, 0, 1)
	IPs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	for _, ip := range IPs {
		if ip.To4() != nil {
			result = append(result, ip)
		}
	}

	return result, nil
}