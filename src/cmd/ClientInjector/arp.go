package main

import (
	"bytes"
	"cmd/ClientInjector/network"
	"log"
	"net"
	"sync/atomic"

	"github.com/google/gopacket/layers"
)

var hwAddrBcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
var hwAddrZero = net.HardwareAddr{0, 0, 0, 0, 0, 0}

type ArpContext struct {
	macAddr net.HardwareAddr
	ipAddr  atomic.Value
	arpIn   chan *layers.ARP
}

func ConstructArpClient(macAddr net.HardwareAddr) (*ArpClient, *ArpContext) {
	c := new(ArpClient)

	c.ctx.macAddr = macAddr
	c.ctx.ipAddr.Store(net.IPv4zero)
	c.ctx.arpIn = make(chan *layers.ARP, 100)
	go c.manageArpPacket()

	return c, &c.ctx
}

type ArpClient struct {
	ctx ArpContext
}

func (self *ArpClient) manageArpPacket() {
	var arpRcv *layers.ARP

	for {
		arpRcv = <-self.ctx.arpIn
		ipAddr := self.ctx.ipAddr.Load().(net.IP)

		if !bytes.Equal(arpRcv.DstHwAddress, hwAddrBcast) && !bytes.Equal(arpRcv.DstHwAddress, self.ctx.macAddr) {
			log.Println(self.ctx.macAddr, "Recieve ARP request for", ipAddr, " but it is ignored because DstMacAddr is inconsistent ")
			continue
		}

		log.Println(self.ctx.macAddr, "Recieve ARP request for", ipAddr)

		eth := &layers.Ethernet{
			SrcMAC:       self.ctx.macAddr,
			DstMAC:       arpRcv.SourceHwAddress,
			EthernetType: layers.EthernetTypeARP,
		}

		arp := &layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPReply,
			SourceHwAddress:   self.ctx.macAddr,
			SourceProtAddress: ipAddr,
			DstHwAddress:      arpRcv.SourceHwAddress,
			DstProtAddress:    arpRcv.SourceProtAddress,
		}

		network.SentPacket(eth, arp)
	}
}

func (self *ArpClient) sendGratuitousARP() error {
	ipAddr := self.ctx.ipAddr.Load().(net.IP)

	eth := &layers.Ethernet{
		SrcMAC:       self.ctx.macAddr,
		DstMAC:       hwAddrBcast,
		EthernetType: layers.EthernetTypeARP,
	}

	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   self.ctx.macAddr,
		SourceProtAddress: ipAddr,
		DstHwAddress:      hwAddrZero,
		DstProtAddress:    ipAddr,
	}

	log.Println(self.ctx.macAddr, "Send Gratuitous ARP", ipAddr)

	return network.SentPacket(eth, arp)
}
