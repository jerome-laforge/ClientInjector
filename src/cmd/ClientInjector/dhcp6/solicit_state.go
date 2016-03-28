package dhcp6

import (
	"cmd/ClientInjector/arp"
	"cmd/ClientInjector/layer"
	"cmd/ClientInjector/network"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/jerome-laforge/dhcp6"
)

type solicitState struct{}

func (self solicitState) do(ctx *dhcp6Context) iState {
	ctx.resetLease()

	// Create ClientServer Message
	packet := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		Options:       make(map[dhcp6.OptionCode][][]byte),
		TransactionID: ctx.xid,
	}

	packet.Options.Add(dhcp6.OptionClientID, dhcp6.NewDUIDLL(1, ctx.MacAddr))
	packet.Options.Add(dhcp6.OptionElapsedTime, dhcp6.ElapsedTime(0))

	// Create RelayMessage and encapsulate ClientServer Message
	relayMsg := &dhcp6.RelayMessage{
		MessageType: dhcp6.MessageTypeRelayForw,
		Options:     make(map[dhcp6.OptionCode][][]byte),
	}

	relayOption := new(dhcp6.RelayMessageOption)
	relayOption.SetClientServerMessage(packet)
	relayMsg.Options.Add(dhcp6.OptionRelayMsg, relayOption)

	interfaceId := dhcp6.RelayMessageOption([]byte(ctx.interfaceID))
	relayMsg.Options.Add(dhcp6.OptionInterfaceID, &interfaceId)

	ba, _ := relayMsg.MarshalBinary()

	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       ctx.MacAddr,
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

	var (
		sleep time.Duration
		tries uint = 1
	)

	for {
		// send Solicit
		for err := network.SentPacket(eth, ip, udp, dhcpv6); err != nil; {
			log.Println(ctx.MacAddr, "SOLICIT: error sending solicit", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// sleep = 2s, 4s, 8s, 16s, 32s, 64s, 64s 64s ...
		if tries < 6 {
			sleep = (1 << tries) * time.Second
		} else {
			sleep = 64 * time.Second
		}

		var (
			payload  []byte
			timeout  time.Duration
			deadline = time.Now().Add(sleep)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(ctx.MacAddr, "SOLICIT: timeout", tries)
				goto TIMEOUT
			case payload = <-ctx.dhcpIn:
				// TODO Manage rcv ADVERTISE
				log.Println("SOLICIT: rcv len =", len(payload))
				return nil
			}
		}
	TIMEOUT:
		tries++
	}
}
