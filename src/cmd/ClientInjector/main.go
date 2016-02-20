package main

import (
	"bytes"
	"dhcpv4/util"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	globalHandle     *pcap.Handle
	dhcpClientsByMac = make(map[uint64]*DhcpClient)
	dhcRelay         = false
	option90         = false
	dhcpContextByIp  DhcpContextByIp
)

type DhcpContextByIp struct {
	mutex sync.RWMutex
	dMap  map[uint32]*dhcpContext
}

func (self *DhcpContextByIp) SetIp(ip uint32, dhcpContext *dhcpContext) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	dhcpContextByIp.dMap[ip] = dhcpContext
}

func (self *DhcpContextByIp) ResetIp(ip uint32) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	delete(dhcpContextByIp.dMap, ip)
}

func (self *DhcpContextByIp) Get(ip uint32) (*dhcpContext, bool) {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	dhcpCLient, ok := dhcpContextByIp.dMap[ip]
	return dhcpCLient, ok
}

func init() {
	dhcpContextByIp.dMap = make(map[uint32]*dhcpContext)
}

func main() {
	var (
		paramIfaceName    = flag.String("eth", "eth0", "Define on which interface the customer will bind")
		paramFirstMacAddr = flag.String("mac", "00:00:14:11:19:77", "First mac address use for the first client (incremented by one for each next clients)")
		paramGiADDR       = flag.String("giaddr", "10.0.0.1", "Use as GiADDR into DHCPv4 header")
		paramNbDhcpClient = flag.Uint("nb_dhcp", 1, "Define number of dhcp client")
		paramLogin        = flag.String("login", "%08d", "Define what is use into option90. fmt.Printf and index of dhcp client with range [0, nb_dhcp[ is used.")
		paramPacing       = flag.Duration("pacing", 100*time.Millisecond, "Define the pacing for launch new dhcp client")
	)
	flag.Parse()

	firstMacAddr, err := net.ParseMAC(*paramFirstMacAddr)
	if err != nil {
		panic(err)
	}
	giaddr, err := util.ConvertIpAddrToUint32(*paramGiADDR)
	if err != nil {
		panic(err)
	}

	if globalHandle, err = getPcapHandleFor(*paramIfaceName); err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Println(err)
		}
	}()

	intFirstMacAddr := util.ConvertMax8byteToUint64([]byte(firstMacAddr))
	macAddr := make([]byte, 8)
	util.ConvertUint64To8byte(intFirstMacAddr+uint64(*paramNbDhcpClient)-1, macAddr)

	log.Println("First : Mac Addr", firstMacAddr, " - login", fmt.Sprintf(*paramLogin, 0))
	log.Println("Last  : Mac Addr", net.HardwareAddr(macAddr[2:]), " - login", fmt.Sprintf(*paramLogin, *paramNbDhcpClient-1))

	rand.Seed(time.Now().UTC().UnixNano())

	// Create each DhcpClient
	for i := uint(0); i < *paramNbDhcpClient; i++ {
		macAddr := make([]byte, 8)
		util.ConvertUint64To8byte(intFirstMacAddr+uint64(i), macAddr)
		macAddr = macAddr[2:]

		var dhcpClient *DhcpClient
		if dhcpClient, err = CreateDhcpClient(macAddr, giaddr, fmt.Sprintf(*paramLogin, i)); err != nil {
			log.Printf("interface %v: %v", *paramIfaceName, err)
			os.Exit(1)
		}

		log.Println("DhcpClient created:", dhcpClient)
		dhcpClientsByMac[intFirstMacAddr+uint64(i)] = dhcpClient

		time.Sleep(*paramPacing)
	}

	// Listen all incoming packets
	go dispatchIncomingPacket()

	// Launch each DhcpClient (DORA and so on)
	for _, dhcpClient := range dhcpClientsByMac {
		dhcpClient.run()
	}

	// Block main goroutine
	select {}
}

func getPcapHandleFor(ifaceName string) (*pcap.Handle, error) {
	// Get interfaces.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, mtu, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	handle.SetBPFFilter(fmt.Sprintf("arp or port %v", bootps))

	return handle, nil

}

func dispatchIncomingPacket() {
	var (
		packetSource = gopacket.NewPacketSource(globalHandle, layers.LayerTypeEthernet)
		in           = packetSource.Packets()
	)

	for {
		packet := <-in
		linkLayer := packet.Layer(layers.LayerTypeEthernet)

		if linkLayer == nil {
			continue
		}

		// udp layer
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			if udpLayer.(*layers.UDP).SrcPort != bootps {
				continue
			}

			appLayer := packet.ApplicationLayer()
			if appLayer == nil {
				continue
			}

			macAddr := util.ConvertMax8byteToUint64([]byte(linkLayer.(*layers.Ethernet).DstMAC))

			if dhcpClient, ok := dhcpClientsByMac[macAddr]; ok {
				dhcpClient.ctx.dhcpIn <- appLayer.Payload()
			}

			// next packet
			continue
		}

		// arp layer
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			arpLayer := layer.(*layers.ARP)
			if arpLayer.Operation != layers.ARPRequest {
				continue
			}

			if !bytes.Equal([]byte(arpLayer.DstHwAddress), []byte(hwAddrBcast)) {
				continue
			}

			if dhcpClient, ok := dhcpContextByIp.Get(util.Convert4byteToUint32(arpLayer.DstProtAddress)); ok {
				dhcpClient.arpIn <- arpLayer
			}

			// next packet
			continue
		}

	}

}

func sentMsg(layers ...gopacket.SerializableLayer) error {
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	if err := gopacket.SerializeLayers(buf, opts, layers...); err != nil {
		return err
	}

	if err := globalHandle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
