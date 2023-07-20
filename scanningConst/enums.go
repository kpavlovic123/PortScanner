package scanningconst

type PortValue string

const (
	OPEN        PortValue = "open"
	CLOSED      PortValue = "closed"
	FILTERED    PortValue = "filtered"
	TCP_TIMEOUT int       = 1500
	UDP_TIMEOUT int       = 1500
)
