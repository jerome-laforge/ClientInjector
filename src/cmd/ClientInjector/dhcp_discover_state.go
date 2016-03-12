package main

import (
	"bytes"
	"cmd/ClientInjector/network"
	"dhcpv4"
	"dhcpv4/option"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"
)

type discoverState struct{}

func (_ discoverState) do(ctx *dhcpContext) iState {
	ctx.resetLease()

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
		SrcPort: network.Bootpc,
		DstPort: network.Bootps,
	}
	udp.SetNetworkLayerForChecksum(ipv4)

	buf := GetBuffer()
	defer ReleaseBuffer(buf)

	discover := new(dhcpv4.DhcpPacket)
	discover.ConstructWithPreAllocatedBuffer(buf, option.DHCPDISCOVER)
	discover.SetMacAddr(ctx.macAddr)
	discover.SetXid(ctx.xid)

	if dhcRelay {
		discover.SetGiAddr(ctx.giaddr)
		discover.AddOption(generateOption82(ctx.macAddr))
	}

	if option90 {
		discover.AddOption(generateOption90(ctx.login))
	}

	bootp := &PayloadLayer{
		contents: discover.Raw,
	}

	var sleep time.Duration
	var tries uint = 1

	for {
		// send discover
		for err := network.SentPacket(eth, ipv4, udp, bootp); err != nil; {
			log.Println(ctx.macAddr, "DISCOVER: error sending discover", err)
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
				log.Println(ctx.macAddr, "DISCOVER: timeout", tries)
				goto TIMEOUT
			case payload = <-ctx.dhcpIn:
				dp, err := dhcpv4.Parse(payload)
				if err != nil {
					// it is not DHCP packet...
					continue
				}

				if !bytes.Equal(ctx.xid, dp.GetXid()) {
					// bug of DHCP Server ?
					log.Println(ctx.macAddr, fmt.Sprintf("DISCOVER: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(ctx.xid), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPOFFER:
						ctx.ipAddr.Store(dp.GetYourIp())
						ctx.serverIp = dp.GetNextServerIp()
						err := ctx.arpClient.sendGratuitousARP()
						if err != nil {
							fmt.Println(ctx.macAddr, "send gratuitousARP error", err)
						}

						ctx.t0, ctx.t1, ctx.t2 = extractAllLeaseTime(dp)

						return requestSelectState{}
					default:
						log.Println(ctx.macAddr, fmt.Sprintf("DISCOVER: Unexcpected message [Excpected: %s] [Actual: %s]", option.DHCPDISCOVER, msgType))
						continue
					}
				} else {
					log.Println(ctx.macAddr, "DISCOVER: Option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
		tries++
	}
}
