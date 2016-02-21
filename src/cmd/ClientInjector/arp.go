package main

import (
	"dhcpv4/util"
	"encoding/binary"
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

type ArpClient interface {
	sendGratuitousARP() error
}

func ConstructArpClient(macAddr net.HardwareAddr) (ArpClient, *ArpContext, error) {
	c := new(arpClient)

	c.ctx.macAddr = macAddr
	c.ctx.ipAddr.Store(uint32(0))
	c.ctx.arpIn = make(chan *layers.ARP, 100)
	go c.manageArpPacket()

	return c, &c.ctx, nil
}

type arpClient struct {
	ctx ArpContext
}

func (self *arpClient) manageArpPacket() {
	var arpRcv *layers.ARP

	for {
		arpRcv = <-self.ctx.arpIn
		ipAddr := self.ctx.ipAddr.Load().(uint32)

		log.Println(self.ctx.macAddr, "Recieve ARP request for", util.ConvertUint32ToIpAddr(ipAddr))

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
			SourceProtAddress: convertUint32ToByte(ipAddr),
			DstHwAddress:      arpRcv.SourceHwAddress,
			DstProtAddress:    arpRcv.SourceProtAddress,
		}

		sentMsg(eth, arp)
	}
}

func (self *arpClient) sendGratuitousARP() error {
	ipAddr := self.ctx.ipAddr.Load().(uint32)

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
		SourceProtAddress: convertUint32ToByte(ipAddr),
		DstHwAddress:      hwAddrZero,
		DstProtAddress:    convertUint32ToByte(ipAddr),
	}

	log.Println(self.ctx.macAddr, "Send Gratuitous ARP", util.ConvertUint32ToIpAddr(ipAddr))

	return sentMsg(eth, arp)
}

func convertUint32ToByte(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
