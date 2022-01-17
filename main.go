package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	var inFlag bool
	flag.BoolVar(&inFlag, "h", false, "Display help")
	flag.Parse()
	if inFlag {
		fmt.Println("\nRun the program with subnet <IP ADDRESS>")
		fmt.Println("Have a nice day...")
		return
	}

	arg := flag.Arg(0)
	if arg == "" {
		log.Fatal("I need a subnet....")
	}

	validIP, validCIDR, maskVal, err := checkCIDR(arg)

	if err != nil {
		log.Fatal("Not a valid IP")
	}

	fmt.Printf("Address: %s\nNetwork: %s\nMask: %s\n\n", validIP, validCIDR, maskVal)

	allHosts, networkIP, bcastIP, err := Hosts(arg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Network: %s\nBroadcast: %s\n", networkIP, bcastIP)
	fmt.Printf("%d hosts in network\n\n", len(allHosts))

}

func checkCIDR(testIP string) (net.IP, *net.IPNet, net.IP, error) {
	ipCheck, cidrCheck, err := net.ParseCIDR(testIP)
	if err != nil {
		return ipCheck, cidrCheck, net.IP(cidrCheck.Mask), err
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
