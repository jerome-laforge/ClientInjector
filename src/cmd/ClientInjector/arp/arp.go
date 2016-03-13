package arp

import (
	"bytes"
	"cmd/ClientInjector/network"
	"log"
	"net"
	"sync/atomic"

	"github.com/google/gopacket/layers"
)

var (
	HwAddrBcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	HwAddrZero  = net.HardwareAddr{0, 0, 0, 0, 0, 0}
)

type ArpContext struct {
	MacAddr net.HardwareAddr
	IpAddr  atomic.Value
	ArpIn   chan *layers.ARP
}

func ConstructArpClient(macAddr net.HardwareAddr) (*ArpClient, *ArpContext) {
	c := new(ArpClient)

	c.ctx.MacAddr = macAddr
	c.ctx.IpAddr.Store(net.IPv4zero)
	c.ctx.ArpIn = make(chan *layers.ARP, 100)
	go c.manageArpPacket()

	return c, &c.ctx
}

type ArpClient struct {
	ctx ArpContext
}

func (self *ArpClient) manageArpPacket() {
	var arpRcv *layers.ARP

	for {
		arpRcv = <-self.ctx.ArpIn
		ipAddr := self.ctx.IpAddr.Load().(net.IP)

		if !bytes.Equal(arpRcv.DstHwAddress, HwAddrBcast) && !bytes.Equal(arpRcv.DstHwAddress, self.ctx.MacAddr) {
			log.Println(self.ctx.MacAddr, "Recieve ARP request for", ipAddr, " but it is ignored because DstMacAddr is inconsistent ")
			continue
		}

		log.Println(self.ctx.MacAddr, "Recieve ARP request for", ipAddr)

		eth := &layers.Ethernet{
			SrcMAC:       self.ctx.MacAddr,
			DstMAC:       arpRcv.SourceHwAddress,
			EthernetType: layers.EthernetTypeARP,
		}

		arp := &layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPReply,
			SourceHwAddress:   self.ctx.MacAddr,
			SourceProtAddress: ipAddr,
			DstHwAddress:      arpRcv.SourceHwAddress,
			DstProtAddress:    arpRcv.SourceProtAddress,
		}

		network.SentPacket(eth, arp)
	}
}

func (self *ArpClient) SendGratuitousARP() error {
	ipAddr := self.ctx.IpAddr.Load().(net.IP)

	eth := &layers.Ethernet{
		SrcMAC:       self.ctx.MacAddr,
		DstMAC:       HwAddrBcast,
		EthernetType: layers.EthernetTypeARP,
	}

	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   self.ctx.MacAddr,
		SourceProtAddress: ipAddr,
		DstHwAddress:      HwAddrZero,
		DstProtAddress:    ipAddr,
	}

	log.Println(self.ctx.MacAddr, "Send Gratuitous ARP", ipAddr)

	return network.SentPacket(eth, arp)
}
