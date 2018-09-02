package arpScan

import (
	"net"
	"time"

	"github.com/mdlayher/arp"
)

func ipCount(cider string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cider)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1:len(ips)-1], nil
}

func inc(ip net.IP) {
	for i := len(ip) -1; i >=0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func ArpScan(network, ifName string, timeout time.Duration) (hwAddrList []string, err error) {
	ifIndex, err := net.InterfaceByName(ifName)
	if err != nil {
		return nil, err
	}

	ipStrings, err := ipCount(network)
	if err != nil {
		return nil, err
	}

	conn, err := arp.Dial(ifIndex)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	for _, ipString := range ipStrings {
		if err := conn.SetDeadline(time.Now().Add(timeout*time.Millisecond)); err != nil {
			continue
		} else {
			targetIp := net.ParseIP(ipString).To4()
			hwAddr, _ := conn.Resolve(targetIp)
			if hwAddr.String() == "" {
				continue
			} else {
				hwAddrList = append(hwAddrList, hwAddr.String())
			}
		}
	}

	return hwAddrList, nil
}