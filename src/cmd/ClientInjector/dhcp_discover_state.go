package main

import (
	"dhcpv4"
	"dhcpv4/option"
	"dhcpv4/util"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"encoding/hex"

	"github.com/google/gopacket/layers"
)

type discoverState struct {
	dhcpContext
}

func (self discoverState) do() iState {
	macAddr := self.macAddr.Load().(net.HardwareAddr)

	self.dhcpContext.resetLease()

	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       macAddr,
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
	discover.SetMacAddr([]byte(macAddr))
	discover.SetXid(self.xid)

	if dhcRelay {
		discover.SetGiAddr(self.giaddr)
		discover.AddOption(generateOption82([]byte(macAddr)))
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
			log.Println(macAddr, "error sending discover", err)
			time.Sleep(2 * time.Second)
			continue
		}

		sleep = time.Duration(math.Min(2*math.Pow(2, float64(retries)), 64)) * time.Second

		var (
			payload  []byte
			timeout  time.Duration
			deadline = time.Now().Add(sleep)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(macAddr, "DISCOVER: timeout", retries)
				goto TIMEOUT
			case payload = <-self.dhcpIn:
				dp, err := dhcpv4.Parse(payload)
				if err != nil {
					// it is not DHCP packet...
					continue
				}

				if self.dhcpContext.xid != util.Convert4byteToUint32(dp.GetXid()) {
					expectedXid := make([]byte, 4)
					util.ConvertUint32To4byte(self.dhcpContext.xid, expectedXid)

					log.Println(macAddr, fmt.Sprintf("DISCOVER: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(expectedXid), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPOFFER:
						self.ipAddr.Store(dp.GetYourIp())
						self.serverIp = dp.GetNextServerIp()
						err := self.arpClient.sendGratuitousARP()
						if err != nil {
							fmt.Println(macAddr, "send gratuitousARP error", err)
						}

						self.t0, self.t1, self.t2 = extractAllLeaseTime(dp)

						return requestSelectState{
							dhcpContext: self.dhcpContext,
						}
					default:
						log.Println(macAddr, fmt.Sprintf("DISCOVER: Unexcpected message [Excpected: %s] [Actual: %s]", option.DHCPDISCOVER, msgType))
						continue
					}
				} else {
					log.Println(macAddr, "DISCOVER: Option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
		retries++
	}
}
