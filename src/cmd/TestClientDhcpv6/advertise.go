package main

import (
	"cmd/ClientInjector/arp"
	"cmd/ClientInjector/layer"
	"cmd/ClientInjector/network"
	"flag"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/mdlayher/dhcp6"
)

var IPv6DHCPv6 = net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, 0, 0x02}

func main() {
	paramIfaceName := flag.String("eth", "eth0", "Define on which interface the customer will bind")
	flag.Parse()

	packet := &dhcp6.Packet{
		MessageType: dhcp6.MessageTypeSolicit,
		Options:     make(map[dhcp6.OptionCode][][]byte),
	}

	rand.Seed(time.Now().UTC().UnixNano())

	r := rand.Uint32()

	for i := range packet.TransactionID {
		packet.TransactionID[i] = byte(r>>uint(i)) & 0xFF
	}

	packet.Options.Add(dhcp6.OptionElapsedTime, dhcp6.ElapsedTime(time.Second))

	ba, err := packet.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("TransactionID: 0x%x\n", packet.TransactionID)

	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x00, 0x13, 0x11, 0x19, 0x77},
		DstMAC:       arp.HwAddrBcast,
		EthernetType: layers.EthernetTypeIPv6,
	}
	ip := &layers.IPv6{
		Version:    6,
		SrcIP:      net.IPv6zero,
		DstIP:      IPv6DHCPv6,
		NextHeader: layers.IPProtocolUDP,
		HopLimit:   1,
	}
	udp := &layers.UDP{
		SrcPort: network.Dhcpv6Client,
		DstPort: network.Dhcpv6Server,
	}

	udp.SetNetworkLayerForChecksum(ip)

	dhcpv6 := &layer.PayloadLayer{
		Contents: ba,
	}

	err = network.OpenPcapHandle(*paramIfaceName, "")
	if err != nil {
		log.Fatal(err)
	}
	defer network.Close()

	err = network.SentPacket(eth, ip, udp, dhcpv6)
	if err != nil {
		log.Fatal(err)
	}
}
