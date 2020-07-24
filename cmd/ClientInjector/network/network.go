package network

import (
	"net"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	pcapHandle   *pcap.Handle
	packetSource *gopacket.PacketSource
	in           chan gopacket.Packet
	once         sync.Once
)

func OpenPcapHandle(ifaceName string, BPFilter string) error {
	// Get interfaces.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return err
	}

	// Open up a pcap handle for packet reads/writes.
	pcapHandle, err = pcap.OpenLive(iface.Name, Mtu, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	if BPFilter != "" {
		pcapHandle.SetBPFFilter(BPFilter)
	}

	return nil

}

func Close() {
	if pcapHandle != nil {
		pcapHandle.Close()
		pcapHandle = nil
	}
}

func SentPacket(layers ...gopacket.SerializableLayer) error {
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	if err := gopacket.SerializeLayers(buf, opts, layers...); err != nil {
		return err
	}

	return pcapHandle.WritePacketData(buf.Bytes())
}

func NextPacket() gopacket.Packet {
	once.Do(func() {
		packetSource = gopacket.NewPacketSource(pcapHandle, layers.LayerTypeEthernet)
		in = packetSource.Packets()
	})

	return <-in
}
