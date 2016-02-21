package main

import (
	"dhcpv4"
	"dhcpv4/option"
	"dhcpv4/util"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"
)

type requestSelectState struct {
	dhcpContext
}

func (self requestSelectState) do() iState {
	ipAddr := self.ipAddr.Load().(uint32)
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

	request := new(dhcpv4.DhcpPacket)
	request.ConstructWithPreAllocatedBuffer(buf, option.DHCPREQUEST)
	request.SetXid(self.xid)
	request.SetMacAddr([]byte(self.macAddr))

	opt50 := new(option.Option50RequestedIpAddress)
	opt50.Construct(ipAddr)
	request.AddOption(opt50)

	opt54 := new(option.Option54DhcpServerIdentifier)
	opt54.Construct(self.serverIp)
	request.AddOption(opt54)

	opt61 := new(option.Option61ClientIdentifier)
	opt61.Construct(byte(1), self.macAddr)
	request.AddOption(opt61)

	if dhcRelay {
		request.SetGiAddr(self.giaddr)
		request.AddOption(generateOption82([]byte(self.macAddr)))
	}

	if option90 {
		request.AddOption(generateOption90(self.login))
	}

	bootp := &PayloadLayer{
		contents: request.Raw,
	}

	for {
		// send request
		for err := sentMsg(eth, ipv4, udp, bootp); err != nil; {
			log.Println(self.macAddr, "SELECT: error sending request", err)
			time.Sleep(2 * time.Second)
		}

		var (
			payload  []byte
			timeout  time.Duration
			deadline = time.Now().Add(2 * time.Second)
		)

		for {
			timeout = deadline.Sub(time.Now())
			select {
			case <-time.After(timeout):
				log.Println(self.macAddr, "SELECT: timeout")
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

					log.Println(self.macAddr, fmt.Sprintf("SELECT: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(expectedXid), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPACK:
						dhcpContextByIp.SetIp(ipAddr, &self.dhcpContext)
						self.t0, self.t1, self.t2 = extractAllLeaseTime(dp)

						return sleepState{
							dhcpContext: self.dhcpContext,
						}
					case option.DHCPNAK:
						log.Println(self.macAddr, "SELECT: receive NAK")
						return discoverState{
							dhcpContext: self.dhcpContext,
						}
					default:
						log.Println(self.macAddr, fmt.Sprintf("SELECT: unexpected message [Excpected: %s] [Actual: %s]", option.DHCPACK, msgType))
						continue
					}
				} else {
					log.Println(self.macAddr, "SELECT: option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
	}
}
