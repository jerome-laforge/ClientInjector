package main

import (
	"cmd/ClientInjector/network"
	"dhcpv4/util"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
)

var (
	dhcpClientsByMac = make(map[uint64]*DhcpClient)
	dhcRelay         bool
	option90         bool
	dhcpContextByIp  DhcpContextByIp
)

type DhcpContextByIp struct {
	mutex sync.RWMutex
	dMap  map[uint32]*dhcpContext
}

func (self *DhcpContextByIp) SetIp(ip uint32, dhcpContext *dhcpContext) {
	self.mutex.Lock()
	dhcpContextByIp.dMap[ip] = dhcpContext
	self.mutex.Unlock()
}

func (self *DhcpContextByIp) ResetIp(ip uint32) {
	self.mutex.Lock()
	delete(dhcpContextByIp.dMap, ip)
	self.mutex.Unlock()
}

func (self *DhcpContextByIp) Get(ip uint32) (*dhcpContext, bool) {
	self.mutex.RLock()
	dhcpCLient, ok := dhcpContextByIp.dMap[ip]
	self.mutex.RUnlock()

	return dhcpCLient, ok
}

func init() {
	dhcpContextByIp.dMap = make(map[uint32]*dhcpContext)
}

func main() {
	var (
		paramIfaceName    = flag.String("eth", "eth0", "Define on which interface the customer will bind")
		paramFirstMacAddr = flag.String("mac", "00:00:13:11:19:77", "First mac address use for the first client (incremented by one for each next clients)")
		paramGiADDR       = flag.String("giaddr", "10.0.0.1", "Use as GiADDR into DHCPv4 header")
		paramNbDhcpClient = flag.Uint("nb_dhcp", 1, "Define number of dhcp client")
		paramLogin        = flag.String("login", "%08d", "Define what is use into option90. fmt.Printf and index of dhcp client with range [0, nb_dhcp[ is used.")
		paramPacing       = flag.Duration("pacing", 100*time.Millisecond, "Define the pacing for launch new dhcp client")
		paramNoOpt90      = flag.Bool("no90", false, "No option 90")
		paramNoRelay      = flag.Bool("noRelay", false, "No relay")
	)
	flag.Parse()

	dhcRelay = !*paramNoRelay
	option90 = !*paramNoOpt90

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

	// Create each DhcpClient
	for i := uint(0); i < *paramNbDhcpClient; i++ {
		macAddr := make([]byte, 8)
		util.ConvertUint64To8byte(intFirstMacAddr+uint64(i), macAddr)
		// Reduce the byte array, as mac addr is only on 6 bytes with BigEndian format
		macAddr = macAddr[2:]

		dhcpClient := CreateDhcpClient(macAddr, giaddr, fmt.Sprintf(*paramLogin, i))

		log.Println("DhcpClient created:", dhcpClient)
		dhcpClientsByMac[intFirstMacAddr+uint64(i)] = dhcpClient
	}

	// Listen all incoming packets
	go dispatchIncomingPacket()

	// Launch each DhcpClient (DORA and so on)
	for _, dhcpClient := range dhcpClientsByMac {
		dhcpClient.run()
		time.Sleep(*paramPacing)
	}

	// Block main goroutine
	select {}
}

func dispatchIncomingPacket() {
	for {
		packet := network.NextPacket()

		// DHCP
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			if udpLayer.(*layers.UDP).SrcPort != network.Bootps {
				continue
			}

			appLayer := packet.ApplicationLayer()
			if appLayer == nil {
				continue
			}

			macAddr := util.ConvertMax8byteToUint64(packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet).DstMAC)

			if dhcpClient, ok := dhcpClientsByMac[macAddr]; ok {
				dhcpClient.ctx.dhcpIn <- appLayer.Payload()
			}

			// next packet
			continue
		}

		// ARP
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			arpLayer := layer.(*layers.ARP)
			if arpLayer.Operation != layers.ARPRequest {
				continue
			}

			if dhcpContext, ok := dhcpContextByIp.Get(util.Convert4byteToUint32(arpLayer.DstProtAddress)); ok {
				dhcpContext.arpIn <- arpLayer
			}

			// next packet
			continue
		}

	}

}
