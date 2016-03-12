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

type requestRenewState struct{}

func (_ requestRenewState) do(ctx *dhcpContext) iState {
	// TODO unicast to self.ServerIp
	ipAddr := ctx.ipAddr.Load().(uint32)
	// Set up all the layers' fields we can.

	eth := &layers.Ethernet{
		SrcMAC:       ctx.macAddr,
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
	request.SetXid(ctx.xid)
	request.SetMacAddr(ctx.macAddr)

	opt50 := new(option.Option50RequestedIpAddress)
	opt50.Construct(ipAddr)
	request.AddOption(opt50)

	opt54 := new(option.Option54DhcpServerIdentifier)
	opt54.Construct(ctx.serverIp)
	request.AddOption(opt54)

	opt61 := new(option.Option61ClientIdentifier)
	opt61.Construct(byte(1), ctx.macAddr)
	request.AddOption(opt61)

	if option90 {
		request.AddOption(generateOption90(ctx.login))
	}

	if dhcRelay && ipAddr == 0 {
		request.SetGiAddr(ctx.giaddr)
		request.AddOption(generateOption82(ctx.macAddr))
	}

	bootp := &PayloadLayer{
		contents: request.Raw,
	}

	for {
		// send request
		for err := sentMsg(eth, ipv4, udp, bootp); err != nil; {
			log.Println(ctx.macAddr, "RENEW: error sending request", err)
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
				log.Println(ctx.macAddr, "RENEW: timeout")

				return timeoutRenewState{}
			case payload = <-ctx.dhcpIn:
				dp, err := dhcpv4.Parse(payload)
				if err != nil {
					// it is not DHCP packet...
					continue
				}

				if !bytes.Equal(ctx.xid, dp.GetXid()) {
					// bug of DHCP Server ?
					log.Println(ctx.macAddr, fmt.Sprintf("RENEW: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(ctx.xid), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPACK:
						ctx.t0, ctx.t1, ctx.t2 = extractAllLeaseTime(dp)
						return sleepState{}
					case option.DHCPNAK:
						log.Println(ctx.macAddr, "RENEW: receive NAK")
						return discoverState{}
					default:
						log.Println(ctx.macAddr, fmt.Sprintf("RENEW: unexpected message [Excpected: %s] [Actual: %s]", option.DHCPACK, msgType))
						continue
					}
				} else {
					log.Println(ctx.macAddr, "RENEW: option 53 is missing")
					continue
				}
			}
		}
	}
}
