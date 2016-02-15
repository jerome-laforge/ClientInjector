package main

import (
	"bytes"
	"dhcpv4"
	"dhcpv4/option"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type discoverState struct {
	dhcpContext
}

func (self discoverState) do() iState {
	in := self.packetsource.Packets()

	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       self.macAddr,
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

	discover := new(dhcpv4.DhcpPacket)
	discover.Construct(option.DHCPDISCOVER)
	discover.SetMacAddr([]byte(self.macAddr))
	discover.SetXid(self.xid)

	bootp := &PayloadLayer{
		contents: discover.Raw,
	}

	var sleep time.Duration
	var retries int

	for {
		// send discover
		for err := sentMsg(self.handle, eth, ipv4, udp, bootp); err != nil; {
			log.Println(self.macAddr, "error sending discover", err)
			time.Sleep(2 * time.Second)
			continue
		}

		sleep = time.Duration(math.Min(2*math.Pow(2, float64(retries)), 64)) * time.Second

		var (
			packet   gopacket.Packet
			timeout  time.Duration
			deadline = time.Now().Add(sleep)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(self.macAddr, "timeout")
				goto TIMEOUT
			case packet = <-in:
				linkLayer := packet.Layer(layers.LayerTypeEthernet)

				// Is it for me?
				if !bytes.Equal([]byte(linkLayer.(*layers.Ethernet).DstMAC), self.macAddr) {
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
					case option.DHCPOFFER:
						self.ipAddr = dp.GetYourIp()
						self.ServerIp = dp.GetNextServerIp()
						err := self.arpClient.sendGratuitousARP()
						if err != nil {
							fmt.Println(self.macAddr, "send gratuitousARP error", err)
						}

						self.t0, self.t1, self.t2 = extractAllLeaseTime(dp)

						return requestSelectState{
							dhcpContext: self.dhcpContext,
						}
					default:
						log.Println(self.macAddr, fmt.Sprintf("Unexcpected message [Excpected: %s] [Actual: %s]", option.DHCPDISCOVER, msgType))
						continue
					}
				} else {
					log.Println(self.macAddr, "Option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
		retries++
	}
}
