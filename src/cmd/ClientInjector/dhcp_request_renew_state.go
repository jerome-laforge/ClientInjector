package main

import (
	"bytes"
	"dhcpv4"
	"dhcpv4/option"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type requestRenewState struct {
	dhcpContext
}

func (self requestRenewState) do() iState {
	// TODO unicast to self.ServerIp
	in := self.packetSource.Packets()
	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       self.macAddr.Load().(net.HardwareAddr),
		DstMAC:       hwAddrBcast,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ipv4 := &layers.IPv4{
		Version:  4,
		TTL:      255,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    net.IPv4zero,
		DstIP:    net.IPv4bcast,
	}
	udp := &layers.UDP{
		SrcPort: 68,
		DstPort: 67,
	}
	udp.SetNetworkLayerForChecksum(ipv4)

	buf := GetBuffer()
	defer ReleaseBuffer(buf)

	request := new(dhcpv4.DhcpPacket)
	request.ConstructWithPreAllocatedBuffer(buf, option.DHCPREQUEST)
	request.SetXid(self.xid)
	request.SetMacAddr([]byte(self.macAddr.Load().(net.HardwareAddr)))

	if self.ipAddr.Load().(uint32) == 0 {
		request.SetGiAddr(self.giaddr)
		request.AddOption(generateOption82([]byte(self.macAddr.Load().(net.HardwareAddr))))
	}

	opt50 := new(option.Option50RequestedIpAddress)
	opt50.Construct(self.ipAddr.Load().(uint32))
	request.AddOption(opt50)

	opt54 := new(option.Option54DhcpServerIdentifier)
	opt54.Construct(self.serverIp)
	request.AddOption(opt54)

	opt61 := new(option.Option61ClientIdentifier)
	opt61.Construct(byte(1), self.macAddr.Load().(net.HardwareAddr))
	request.AddOption(opt61)

	request.AddOption(generateOption90(self.login))

	bootp := &PayloadLayer{
		contents: request.Raw,
	}

	for {
		// send request
		for err := sentMsg(self.handle, eth, ipv4, udp, bootp); err != nil; {
			log.Println(self.macAddr, "error sending request", err)
			time.Sleep(2 * time.Second)
		}

		var (
			packet   gopacket.Packet
			timeout  time.Duration
			deadline = time.Now().Add(2 * time.Second)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(self.macAddr, "timeout")

				return timeoutRenewState{
					dhcpContext: self.dhcpContext,
				}
			case packet = <-in:
				linkLayer := packet.Layer(layers.LayerTypeEthernet)

				// Is it for me?
				if !bytes.Equal([]byte(linkLayer.(*layers.Ethernet).DstMAC), self.macAddr.Load().(net.HardwareAddr)) {
					// no, ignore this packet.
					continue
				}

				appLayer := packet.ApplicationLayer()
				if appLayer == nil {
					continue
				}

				dp, err := dhcpv4.Parse(appLayer.Payload())
				if err != nil {
					// it is not DHCP packet...
					continue
				}
				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPACK:
						self.t0, self.t1, self.t2 = extractAllLeaseTime(dp)

						return sleepState{
							dhcpContext: self.dhcpContext,
						}
					case option.DHCPNAK:
						log.Println(self.macAddr, "Receive NAK")
						return discoverState{
							dhcpContext: self.dhcpContext,
						}
					default:
						log.Println(self.macAddr, fmt.Sprintf("Unexpected message [Excpected: %s] [Actual: %s]", option.DHCPACK, msgType))
						continue
					}
				} else {
					log.Println(self.macAddr, "Option 53 is missing")
					continue
				}
			}
		}
	}
}
