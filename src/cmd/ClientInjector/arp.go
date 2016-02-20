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
	macAddr atomic.Value
	ipAddr  atomic.Value
	arpIn   chan *layers.ARP
}

type ArpClient interface {
	sendGratuitousARP() error
}

func ConstructArpClient(macAddr net.HardwareAddr) (ArpClient, *ArpContext, error) {
	c := new(arpClient)

	c.ctx.macAddr.Store(macAddr)
	c.ctx.ipAddr.Store(uint32(0))
	c.ctx.arpIn = make(chan *layers.ARP, 100)
	go c.manageArpPacket()

	return c, &c.ctx, nil
}

type arpClient struct {
	ctx ArpContext
}

func (self *arpClient) manageArpPacket() {
	var (
		macAddr = self.ctx.macAddr.Load().(net.HardwareAddr)
		arpRcv  *layers.ARP
	)

	for {
		arpRcv = <-self.ctx.arpIn
		ipAddr := self.ctx.ipAddr.Load().(uint32)

		log.Println(macAddr, "Recieve ARP request for", util.ConvertUint32ToIpAddr(ipAddr))

		eth := &layers.Ethernet{
			SrcMAC:       macAddr,
			DstMAC:       arpRcv.SourceHwAddress,
			EthernetType: layers.EthernetTypeARP,
		}

		arp := &layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPReply,
			SourceHwAddress:   []byte(macAddr),
			SourceProtAddress: convertUint32ToByte(ipAddr),
			DstHwAddress:      arpRcv.SourceHwAddress,
			DstProtAddress:    arpRcv.SourceProtAddress,
		}

		sentMsg(eth, arp)
	}
}

func (self *arpClient) sendGratuitousARP() error {
	var (
		macAddr = self.ctx.macAddr.Load().(net.HardwareAddr)
		ipAddr  = self.ctx.ipAddr.Load().(uint32)
	)

	eth := &layers.Ethernet{
		SrcMAC:       macAddr,
		DstMAC:       hwAddrBcast,
		EthernetType: layers.EthernetTypeARP,
	}

	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   macAddr,
		SourceProtAddress: convertUint32ToByte(ipAddr),
		DstHwAddress:      hwAddrZero,
		DstProtAddress:    convertUint32ToByte(ipAddr),
	}

	log.Println(macAddr, "Send Gratuitous ARP", util.ConvertUint32ToIpAddr(ipAddr))

	return sentMsg(eth, arp)
}

func convertUint32ToByte(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
