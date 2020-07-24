package dhcp4

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket/layers"

	"github.com/jerome-laforge/ClientInjector/cmd/ClientInjector/arp"
	"github.com/jerome-laforge/ClientInjector/cmd/ClientInjector/layer"
	"github.com/jerome-laforge/ClientInjector/cmd/ClientInjector/network"
	"github.com/jerome-laforge/ClientInjector/dhcpv4"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/option"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/util"
)

type requestSelectState struct{}

func (_ requestSelectState) do(ctx *dhcpContext) iState {
	ipAddr := ctx.IpAddr.Load().(net.IP)

	// Set up all the layers' fields we can.
	eth := &layers.Ethernet{
		SrcMAC:       ctx.MacAddr,
		DstMAC:       arp.HwAddrBcast,
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

	buf := network.GetBuffer()
	defer network.ReleaseBuffer(buf)

	request := new(dhcpv4.DhcpPacket)
	request.ConstructWithPreAllocatedBuffer(buf, option.DHCPREQUEST)
	request.SetXid(ctx.xid[:])
	request.SetMacAddr(ctx.MacAddr)

	opt50 := new(option.Option50RequestedIpAddress)
	opt50.Construct(util.Convert4byteToUint32(ipAddr))
	request.AddOption(opt50)

	opt54 := new(option.Option54DhcpServerIdentifier)
	opt54.Construct(ctx.serverIp)
	request.AddOption(opt54)

	opt61 := new(option.Option61ClientIdentifier)
	opt61.Construct(byte(1), ctx.MacAddr)
	request.AddOption(opt61)

	if DhcRelay {
		request.SetGiAddr(ctx.giaddr)
		request.AddOption(generateOption82(ctx.MacAddr))
	}

	if Option90 {
		request.AddOption(generateOption90(ctx.login))
	}

	bootp := &layer.PayloadLayer{
		Contents: request.Raw,
		LT:       layers.LayerTypeDHCPv4,
	}

	for {
		// send request
		for err := network.SentPacket(eth, ipv4, udp, bootp); err != nil; {
			log.Println(ctx.MacAddr, "SELECT: error sending request", err)
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
				log.Println(ctx.MacAddr, "SELECT: timeout")
				goto TIMEOUT
			case payload = <-ctx.dhcpIn:
				dp, err := dhcpv4.Parse(payload)
				if err != nil {
					// it is not DHCP packet...
					continue
				}

				if !bytes.Equal(ctx.xid[:], dp.GetXid()) {
					// bug of DHCP Server ?
					log.Println(ctx.MacAddr, fmt.Sprintf("SELECT: unexpected xid [Expected: 0x%v] [Actual: 0x%v]", hex.EncodeToString(ctx.xid[:]), hex.EncodeToString(dp.GetXid())))
					continue
				}

				if msgType, err := dp.GetTypeMessage(); err == nil {
					switch msgType {
					case option.DHCPACK:
						arp.MapArpByIp.Set(ipAddr, ctx.ArpIn)
						ctx.t0, ctx.t1, ctx.t2 = extractAllLeaseTime(dp)
						return sleepState{}
					case option.DHCPNAK:
						log.Println(ctx.MacAddr, "SELECT: receive NAK")
						return discoverState{}
					default:
						log.Println(ctx.MacAddr, fmt.Sprintf("SELECT: unexpected message [Excpected: %s] [Actual: %s]", option.DHCPACK, msgType))
						continue
					}
				} else {
					log.Println(ctx.MacAddr, "SELECT: option 53 is missing")
					continue
				}
			}
		}
	TIMEOUT:
	}
}
