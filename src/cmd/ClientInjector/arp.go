package main

import (
	"bytes"
	"dhcpv4/util"
	"encoding/binary"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var hwAddrBcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
var hwAddrZero = net.HardwareAddr{0, 0, 0, 0, 0, 0}

type ArpContext struct {
	macAddr net.HardwareAddr
	ipAddr  uint32
}

type ArpClient interface {
	sendGratuitousARP() error
}

func ConstructArpClient(ifaceName string, macAddr net.HardwareAddr) (ArpClient, *ArpContext, error) {
	c := new(arpClient)
	h, err := getPcapHandleFor(ifaceName)
	if err != nil {
		return nil, nil, err
	}

	c.arpHandle = h
	c.ctx.macAddr = macAddr
	go c.manageArpPacket()

	return c, &c.ctx, nil
}

type arpClient struct {
	arpHandle *pcap.Handle
	ctx       ArpContext
}

func (self *arpClient) manageArpPacket() {
	src := gopacket.NewPacketSource(self.arpHandle, layers.LayerTypeEthernet)
	in := src.Packets()

	var packet gopacket.Packet
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	for {
		packet = <-in
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil {
			continue
		}
		arpRcv := arpLayer.(*layers.ARP)
		// if arpRcv.Operation != layers.ARPRequest || !bytes.Equal(arpRcv.DstHwAddress, self.clientMAC) {
		if arpRcv.Operation != layers.ARPRequest || util.Convert4byteToUint32(arpRcv.DstProtAddress) != self.ctx.ipAddr {
			if arpRcv.Operation != RARPRequest || !bytes.Equal(arpRcv.DstHwAddress, []byte(self.ctx.macAddr)) {
				continue
			} else {
				log.Println(self.ctx.macAddr, "Recieve RARP request")
			}
		} else {
			log.Println(self.ctx.macAddr, "Recieve ARP request")
		}

		eth := &layers.Ethernet{
			SrcMAC:       self.ctx.macAddr,
			DstMAC:       arpRcv.SourceHwAddress,
			EthernetType: layers.EthernetTypeARP,
		}

		var op uint16
		if arpRcv.Operation == layers.ARPRequest {
			op = layers.ARPReply
		} else { // arpRcv.Operation == RARPRequest
			op = RARPReply
		}
		arp := &layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         op,
			SourceHwAddress:   []byte(self.ctx.macAddr),
			SourceProtAddress: convertUint32ToByte(self.ctx.ipAddr),
			DstHwAddress:      arpRcv.SourceHwAddress,
			DstProtAddress:    arpRcv.SourceProtAddress,
		}

		gopacket.SerializeLayers(buf, opts, eth, arp)
		if err := self.arpHandle.WritePacketData(buf.Bytes()); err != nil {
			log.Println(self.ctx.macAddr, "ARP reply error", err)
			continue
		}
	}
}

func (self *arpClient) sendGratuitousARP() error {
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
		SourceHwAddress:   []byte(self.ctx.macAddr),
		SourceProtAddress: convertUint32ToByte(self.ctx.ipAddr),
		DstHwAddress:      hwAddrZero,
		DstProtAddress:    convertUint32ToByte(self.ctx.ipAddr),
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	gopacket.SerializeLayers(buf, opts, eth, arp)
	log.Println(self.ctx.macAddr, "Send Gratuitous ARP", util.ConvertUint32ToIpAddr(self.ctx.ipAddr))
	return self.arpHandle.WritePacketData(buf.Bytes())
}

func convertUint32ToByte(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

const (
	RARPRequest = 3
	RARPReply   = 4
)
