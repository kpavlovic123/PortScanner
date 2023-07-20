package tcp

import (
	"net"
	sc "portScanner/scanningConst"
	"strconv"
	"strings"
	"time"
)

func Scan(ipAddr net.IPAddr, port string) sc.PortValue {
	p, _ := strconv.Atoi(port)
	tcpAddr := net.TCPAddr{IP: ipAddr.IP, Port: p}
	conn, err := net.DialTimeout("tcp4", tcpAddr.IP.String()+":"+port, time.Millisecond*time.Duration(sc.TCP_TIMEOUT))
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return sc.FILTERED
		} else if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(time.Millisecond * time.Duration(sc.TCP_TIMEOUT))
			return Scan(ipAddr, port)
		} else {
			return sc.CLOSED
		}
	}
	conn.Close()
	return sc.OPEN
}
