package arp

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/google/gopacket/layers"

	"github.com/jerome-laforge/ClientInjector/cmd/ClientInjector/network"
)

var (
	HwAddrBcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	HwAddrZero  = net.HardwareAddr{0, 0, 0, 0, 0, 0}
	MapArpByIp  ArpInByIp
)

type ArpInByIp struct {
	mutex    sync.RWMutex
	arpInMap map[string]chan *layers.ARP
}

func (self *ArpInByIp) Set(ip net.IP, inArp chan *layers.ARP) {
	self.mutex.Lock()
	self.arpInMap[string(ip)] = inArp
	self.mutex.Unlock()
}

func (self *ArpInByIp) Reset(ip net.IP) {
	self.mutex.Lock()
	delete(self.arpInMap, string(ip))
	self.mutex.Unlock()
}

func (self *ArpInByIp) Lookup(ip net.IP) (chan *layers.ARP, bool) {
	self.mutex.RLock()
	inArp, ok := self.arpInMap[string(ip)]
	self.mutex.RUnlock()

	return inArp, ok
}

func init() {
	MapArpByIp.arpInMap = make(map[string]chan *layers.ARP)
}

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

		if !bytes.Equal(arpRcv.DstHwAddress, self.ctx.MacAddr) && !bytes.Equal(arpRcv.DstHwAddress, HwAddrZero) && !bytes.Equal(arpRcv.DstHwAddress, HwAddrBcast) {
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
			Protocol:          getProtocolFor(ipAddr),
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
		Protocol:          getProtocolFor(ipAddr),
		Operation:         layers.ARPRequest,
		SourceHwAddress:   self.ctx.MacAddr,
		SourceProtAddress: ipAddr,
		DstHwAddress:      HwAddrZero,
		DstProtAddress:    ipAddr,
	}

	log.Println(self.ctx.MacAddr, "Send Gratuitous ARP", ipAddr)

	return network.SentPacket(eth, arp)
}

func getProtocolFor(ipAddr net.IP) layers.EthernetType {
	switch len(ipAddr) {
	case net.IPv4len:
		return layers.EthernetTypeIPv4
	case net.IPv6len:
		return layers.EthernetTypeIPv6
	default:
		panic("Unknown ipAddr " + strconv.Itoa(len(ipAddr)))
	}
}
