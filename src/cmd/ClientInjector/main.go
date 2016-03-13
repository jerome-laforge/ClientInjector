package main

import (
	"cmd/ClientInjector/arp"
	"cmd/ClientInjector/dhcp4"
	"cmd/ClientInjector/network"
	"dhcpv4/util"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/google/gopacket/layers"
)

var MapClientsByMac = make(map[uint64]chan []byte)

func main() {
	var (
		paramIfaceName    = flag.String("eth", "eth0", "Define on which interface the customer will bind")
		paramFirstMacAddr = flag.String("mac", "00:00:13:11:19:77", "First mac address use for the first client (incremented by one for each next clients)")
		paramGiADDR       = flag.String("giaddr", "10.0.0.1", "Use as GiADDR into DHCPv4 header")
		paramNbDhcpClient = flag.Uint("nb_dhcp", 1, "Define number of dhcp client")
		paramLogin        = flag.String("login", "%08d", "Define what is use into option90. fmt.Printf and index of dhcp client with range [0, nb_dhcp[ is used.")
		paramPacing       = flag.Duration("pacing", 100*time.Millisecond, "Define the pacing for launch new dhcp client")
		paramNoLogin      = flag.Bool("noLogin", false, "No login (DHCPv4: option 90)")
		paramNoRelay      = flag.Bool("noRelay", false, "No relay")
	)
	flag.Parse()

	dhcp4.DhcRelay = !*paramNoRelay
	dhcp4.Option90 = !*paramNoLogin

	firstMacAddr, err := net.ParseMAC(*paramFirstMacAddr)
	if err != nil {
		log.Fatal(err)
	}
	giaddr, err := util.ConvertIpAddrToUint32(*paramGiADDR)
	if err != nil {
		log.Fatal(err)
	}

	if err = network.OpenPcapHandle(*paramIfaceName); err != nil {
		log.Fatal(err)
	}
	defer network.Close()

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Println(err)
		}
	}()

	intFirstMacAddr := util.ConvertMax8byteToUint64(firstMacAddr)
	macAddr := make([]byte, 8)
	util.ConvertUint64To8byte(intFirstMacAddr+uint64(*paramNbDhcpClient)-1, macAddr)

	log.Println("First : Mac Addr", firstMacAddr, "- login", fmt.Sprintf(*paramLogin, 0))
	log.Println("Last  : Mac Addr", net.HardwareAddr(macAddr[2:]), "- login", fmt.Sprintf(*paramLogin, *paramNbDhcpClient-1))

	rand.Seed(time.Now().UTC().UnixNano())

	{
		var listRunnable = make([]runnable, 0, *paramNbDhcpClient)
		// Create each DhcpClient
		for i := 0; i < int(*paramNbDhcpClient); i++ {
			macAddr := make([]byte, 8)
			util.ConvertUint64To8byte(intFirstMacAddr+uint64(i), macAddr)
			// Reduce the byte array, as mac addr is only on 6 bytes with BigEndian format
			macAddr = macAddr[2:]

			dhcpClient, dhcpIn := dhcp4.CreateClient(macAddr, giaddr, fmt.Sprintf(*paramLogin, i))

			log.Println("DhcpClient created:", dhcpClient)
			MapClientsByMac[intFirstMacAddr+uint64(i)] = dhcpIn
			listRunnable = append(listRunnable, dhcpClient)
		}

		// Listen all incoming packets
		go dispatchIncomingPacket()

		// Launch each DhcpClient (DORA and so on)
		for i := range listRunnable {
			listRunnable[i].Run()
			time.Sleep(*paramPacing)
		}
	}

	// Block main goroutine
	select {}
}

func dispatchIncomingPacket() {
	for {
		packet := network.NextPacket()

		// DHCP
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			// DHCPv4
			if udpLayer.(*layers.UDP).SrcPort == network.Bootps {
				appLayer := packet.ApplicationLayer()
				if appLayer == nil {
					continue
				}

				macAddr := util.ConvertMax8byteToUint64(packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet).DstMAC)

				if dhcpIn, ok := MapClientsByMac[macAddr]; ok {
					dhcpIn <- appLayer.Payload()
				}

				// next packet
				continue
			}

			if udpLayer.(*layers.UDP).SrcPort == network.Dhcpv6Server {
				// TODO DHCPv6
			}
		}

		// ARP
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			arpLayer := layer.(*layers.ARP)
			if arpLayer.Operation != layers.ARPRequest {
				continue
			}

			if arpIn, ok := arp.MapArpByIp.Lookup(arpLayer.DstProtAddress); ok {
				arpIn <- arpLayer
			}

			// next packet
			continue
		}
	}
}

type runnable interface {
	Run()
}
