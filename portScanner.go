package main

import (
	"fmt"
	"net"
	"os"
	sc "portScanner/scanningConst"
	"portScanner/tcp"
	"portScanner/udp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type port struct {
	port             string
	portType         string
	portAvailability sc.PortValue
}

type scanRange struct {
	ports  []string
	ipAddr *net.IPAddr
}

type scanResults struct {
	cp []port
	op []port
	fp []port
}

type sMutex struct {
	mu1, mu2, mu3 sync.Mutex
}

var srn scanRange
var sres scanResults
var mu sMutex
var tcpEnabled, udpEnabled, fileOutput, setsOutput bool
var numOfPorts int

// Čeka portove, šalje paket i razvrstava u skupove
func portWorker(ports <-chan port, wg *sync.WaitGroup) {
	defer wg.Done()
	cp, op, fp := make([]port, 100), make([]port, 100), make([]port, 100)
	i1, i2, i3 := 0, 0, 0
	l1, l2, l3 := len(cp), len(op), len(fp)
	for p := range ports {
		var r sc.PortValue
		switch p.portType {
		case "tcp":
			r = tcp.Scan(*srn.ipAddr, p.port)
		case "udp":
			r = udp.Scan(*srn.ipAddr, p.port)
		}
		p.portAvailability = r

		add := func(i, l int, p port, sp []port) (int, int, []port) {
			if i == l {
				sp = append(sp, p)
				l = len(sp)
				i++
			} else {
				sp[i] = p
				i++
			}
			return i, l, sp
		}

		switch r {
		case sc.CLOSED:
			i1, l1, cp = add(i1, l1, p, cp)
		case sc.OPEN:
			i2, l2, op = add(i2, l2, p, op)
		case sc.FILTERED:
			i3, l3, fp = add(i3, l3, p, fp)

		}
	}

	cp, op, fp = cp[:i1], op[:i2], fp[:i3]

	lockAndAppend := func(res *[]port, a []port, mu *sync.Mutex) {
		mu.Lock()
		*res = append((*res), a...)
		mu.Unlock()
	}

	lockAndAppend(&sres.cp, cp, &mu.mu1)
	lockAndAppend(&sres.op, op, &mu.mu2)
	lockAndAppend(&sres.fp, fp, &mu.mu3)
}

// Pokreće workere, parsira portove i šalje ih workerima
func PortScanner() {
	var wg1 sync.WaitGroup
	ports := make(chan port, numOfPorts)

	numOfWorkers := numOfPorts/10 + 10

	fmt.Println("Starting scan...")
	startTimer := time.Now()

	for i := int64(0); i < int64(numOfWorkers); i++ {
		go portWorker(ports, &wg1)
		wg1.Add(1)
	}

	for _, v := range srn.ports {
		if strings.Contains(v, "-") {
			r := strings.Split(v, "-")
			begin, _ := strconv.Atoi(r[0])
			end, _ := strconv.Atoi(r[1])
			for i := begin; i <= end; i++ {
				if tcpEnabled {
					ports <- port{port: strconv.Itoa(i), portType: "tcp"}
				}
				if udpEnabled {
					ports <- port{port: strconv.Itoa(i), portType: "udp"}
				}
			}
		} else {
			if tcpEnabled {
				ports <- port{port: v, portType: "tcp"}
			}
			if udpEnabled {
				ports <- port{port: v, portType: "udp"}
			}
		}

	}
	close(ports)

	wg1.Wait()

	endTimer := time.Now()

	fmt.Println("Elapsed time:", endTimer.Sub(startTimer))
}

func sortPortArray(x *[]port) {
	sort.Slice(*x, func(i, j int) bool {
		v1, _ := strconv.Atoi((*x)[i].port)
		v2, _ := strconv.Atoi((*x)[j].port)
		return v1 < v2
	})
}

func printPortArray(pa []port) {
	for _, p := range pa {
		if p.portAvailability == sc.FILTERED && p.portType == "udp" {
			fmt.Printf("(%s/%s) is open|filtered\n", p.portType, p.port)
		} else {
			fmt.Printf("(%s/%s) is %s\n", p.portType, p.port, p.portAvailability)
		}
	}
}

func parseArguments() {
	if len(os.Args) == 1 {
		fmt.Println("Destination address required.")
		os.Exit(1)
	}
	args := os.Args[1:]
	var ptc []string
	address := args[len(args)-1]
	for i, v := range args {
		switch v {
		case "-p":
			ptc = strings.Split(args[i+1], ",")
		case "-p-":
			ptc = []string{"1-65535"}
		case "-sT":
			tcpEnabled = true
		case "-sU":
			udpEnabled = true
		case "-o":
			fileOutput = true
		case "-s":
			setsOutput = true
		}
	}

	if ptc == nil {
		ptc = []string{"1-1023"}
	}

	var err error

	srn.ipAddr, err = net.ResolveIPAddr("ip4", address)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	srn.ports = ptc

	if !udpEnabled && !tcpEnabled {
		udpEnabled, tcpEnabled = true, true
	}

	for _, v := range srn.ports {
		if strings.Contains(v, "-") {
			r := strings.Split(v, "-")
			v1, _ := strconv.Atoi(r[0])
			v2, _ := strconv.Atoi(r[1])
			numOfPorts += v2 - v1 + 1
		} else {
			numOfPorts++
		}
	}

	if udpEnabled && tcpEnabled {
		numOfPorts *= 2
	}

	fmt.Println("Number of ports:", numOfPorts)
}

func main() {
	parseArguments()
	PortScanner()
	if fileOutput {
		filepath := "portScannerResults.txt"

		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		fmt.Println("Writing results to file.")

		defer file.Close()

		file.WriteString("OPEN: ")
		for _, v := range sres.op {
			file.WriteString(v.portType + "," + v.port + " ")
		}
		file.WriteString("\nCLOSED: ")
		for _, v := range sres.cp {
			file.WriteString(v.portType + "," + v.port + " ")
		}
		file.WriteString("\nFILTERED: ")
		for _, v := range sres.fp {
			file.WriteString(v.portType + "," + v.port + " ")
		}
		fmt.Println("Wrote files to ", filepath)
	}

	if setsOutput {
		fmt.Print("OPEN: ")
		for _, v := range sres.op {
			fmt.Print(v.portType + "," + v.port + " ")
		}
		fmt.Println()

		fmt.Print("CLOSED: ")
		for _, v := range sres.cp {
			fmt.Print(v.portType + "," + v.port + " ")
		}
		fmt.Println()

		fmt.Print("FILTERED: ")
		for _, v := range sres.fp {
			fmt.Print(v.portType + "," + v.port + " ")
		}
		fmt.Println()
	} else {
		fmt.Println("\nOpen ports:", len(sres.op))
		if len(sres.op) > 0 {
			sortPortArray(&sres.op)
			printPortArray(sres.op)
		}

		fmt.Println("\nFiltered ports:", len(sres.fp))
		if len(sres.fp) < 15 && len(sres.fp) > 0 {
			sortPortArray(&sres.fp)
			printPortArray(sres.fp)
		}

		fmt.Println("\nClosed ports:", len(sres.cp))
		if len(sres.cp) < 15 && len(sres.cp) > 0 {
			sortPortArray(&sres.cp)
			printPortArray(sres.cp)
		}
	}
}
