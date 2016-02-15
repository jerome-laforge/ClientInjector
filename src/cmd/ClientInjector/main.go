package main

import (
	"dhcpv4/util"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/google/gopacket/pcap"
)

func main() {
	ifaceName := flag.String("eth", "eth0", "Define on which interface the customer will bind")
	nbDhcpClient := flag.Uint("nb_dhcp", 1, "Define number of dhcp client")
	pacing := flag.Duration("pacing", 100*time.Millisecond, "Define the pacing for launch new dhcp client")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	firstMacAddr, _ := net.ParseMAC("00:00:14:11:19:77")
	intFirstMacAddr := util.ConvertMax8byteToUint64([]byte(firstMacAddr))

	for i := uint(0); i < *nbDhcpClient; i++ {
		macAddr := make([]byte, 8)
		util.ConvertUint64To8byte(intFirstMacAddr+uint64(i), macAddr)
		macAddr = macAddr[2:]
		if _, err := CreateDhcpClient(*ifaceName, macAddr); err != nil {
			log.Printf("interface %v: %v", *ifaceName, err)
			os.Exit(1)
		}

		time.Sleep(*pacing)
	}

	select {}
}

func getPcapHandleFor(ifaceName string) (*pcap.Handle, error) {
	// Get interfaces.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	// We just look for IPv4 addresses, so try to find if the interface has one.
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return nil, err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return nil, fmt.Errorf("no good IP network found")
	} else if addr.IP[0] == 127 {
		return nil, fmt.Errorf("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return nil, fmt.Errorf("mask means network is too large")
	}

	// Open up a pcap handle for packet reads/writes.
	return pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)

}
