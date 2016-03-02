package main

import (
	"bytes"
	"dhcpv4"
	"dhcpv4/option"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"
)

type discoverState struct {
	dhcpContext
}

func (self discoverState) do() iState {
	self.dhcpContext.resetLease()

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
		SrcPort: bootpc,
		DstPort: bootps,
	}
	udp.SetNetworkLayerForChecksum(ipv4)

	buf := GetBuffer()
	defer ReleaseBuffer(buf)

	discover := new(dhcpv4.DhcpPacket)
	discover.ConstructWithPreAllocatedBuffer(buf, option.DHCPDISCOVER)
	discover.SetMacAddr(self.macAddr)
	discover.SetXid(self.xid)

	if dhcRelay {
		discover.SetGiAddr(self.giaddr)
		discover.AddOption(generateOption82(self.macAddr))
	}

	if option90 {
		discover.AddOption(generateOption90(self.login))
	}

	bootp := &PayloadLayer{
		contents: discover.Raw,
	}

	var sleep time.Duration
	var retries int

	for {
		// send discover
		for err := sentMsg(eth, ipv4, udp, bootp); err != nil; {
			log.Println(self.macAddr, "DISCOVER: error sending discover", err)
			time.Sleep(2 * time.Second)
			continue
		}

		sleep = time.Duration(Min(2*Pow(2, retries), 64)) * time.Second

		var (
			payload  []byte
			timeout  time.Duration
			deadline = time.Now().Add(sleep)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(self.macAddr, "DISCOVER: timeout", retries)
				goto TIMEOUT
			case payload = <-self.dhcpIn:
				dp, err := dhcpv4.Parse(payload)
				if err != nil {
					// it is not DHCP packet...
					continue
				}

				if !bytes.Equal(self.xid, dp.GetXid()) {
					// bug of DHCP Server ?
					log.Println(self.macAddr, fmt.Sprintf("DISCOVER: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(self.xid), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPOFFER:
						self.ipAddr.Store(dp.GetYourIp())
						self.serverIp = dp.GetNextServerIp()
						err := self.arpClient.sendGratuitousARP()
						if err != nil {
							fmt.Println(self.macAddr, "send gratuitousARP error", err)
						}

						self.t0, self.t1, self.t2 = extractAllLeaseTime(dp)

						return requestSelectState{
							dhcpContext: self.dhcpContext,
						}
					default:
						log.Println(self.macAddr, fmt.Sprintf("DISCOVER: Unexcpected message [Excpected: %s] [Actual: %s]", option.DHCPDISCOVER, msgType))
						continue
					}
				} else {
					log.Println(self.macAddr, "DISCOVER: Option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
		retries++
	}
}
