package udp

import (
	"net"
	sc "portScanner/scanningConst"
	"strconv"
	"strings"
	"time"
)

const portScannerMsg string = "PortScannerv1.0"

func Scan(ipAddr net.IPAddr, port string) sc.PortValue {
	p, _ := strconv.Atoi(port)
	udpAddr := net.UDPAddr{IP: ipAddr.IP, Port: p}
	conn, err := net.DialUDP("udp4", nil, &udpAddr)
	if err != nil {
		if strings.Contains(err.Error(), "resource temporarily unavailable") {
			time.Sleep(time.Millisecond * time.Duration(sc.UDP_TIMEOUT))
			return Scan(ipAddr, port)
		} else {
			return sc.CLOSED
		}
	}
	defer conn.Close()
	_, err = conn.Write([]byte(portScannerMsg))
	if err != nil {
		return sc.CLOSED
	}
	conn.SetReadDeadline(time.Now().Add(time.Duration(sc.UDP_TIMEOUT) * time.Millisecond))
	buf := make([]byte, 15)
	n, err := conn.Read(buf)
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return sc.FILTERED
		} else {
			return sc.CLOSED
		}
	}

	if string(buf[:n]) == portScannerMsg {
		return sc.CLOSED
	}

	return sc.OPEN
}
