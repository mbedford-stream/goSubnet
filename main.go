package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fatih/color"
)

func main() {
	var inFlag bool
	var printList bool
	flag.BoolVar(&inFlag, "h", false, "Display help")
	flag.BoolVar(&printList, "p", false, "Print list of IPs in subnet")
	flag.Parse()
	if inFlag {
		fmt.Println("\nRun the program with subnet (-p) <IP ADDRESS/CIDR>")
		fmt.Println("Have a nice day...")
		return
	}

	arg := flag.Arg(0)
	if arg == "" {
		log.Fatal("I need a subnet....")
	}

	validIP, validCIDR, maskVal, err := checkCIDR(arg)

	if maskVal.String() == "255.255.255.255" {
		color.Red("I'm not able to do subnet things with a single IP")
		os.Exit(0)
	}

	if err != nil {
		// fmt.Println(err)
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(0)
	}

	fmt.Printf("Address: %s\nNetwork: %s\nMask: %s\n\n", validIP, validCIDR, maskVal)

	allHosts, networkIP, bcastIP, err := Hosts(arg)
	if err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(0)
	}

	fmt.Printf("Network: %s\nBroadcast: %s\n", networkIP, bcastIP)
	fmt.Printf("%d hosts in network\n\n", len(allHosts))

	if printList {
		for _, i := range allHosts {
			fmt.Println(i)
		}
	}

}

func checkCIDR(testIP string) (net.IP, *net.IPNet, net.IP, error) {
	ipCheck, cidrCheck, err := net.ParseCIDR(testIP)
	if err != nil {
		// return ipCheck, cidrCheck, net.IP(cidrCheck.Mask), err
		return nil, nil, nil, err
	}

	return ipCheck, cidrCheck, net.IP(cidrCheck.Mask), nil

}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Hosts(cidr string) ([]string, string, string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, "", "", err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// return separate values for host IPs, network address, and broadcast address
	return ips[1 : len(ips)-1], ips[0], ips[len(ips)-1], nil
}
