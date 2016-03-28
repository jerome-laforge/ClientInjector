package dhcp6

import (
	"cmd/ClientInjector/arp"
	"net"
	"time"
)

var IPv6DHCPv6 = net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, 0, 0x02}

type iState interface {
	do(*dhcp6Context) iState
}

type dhcp6Context struct {
	*arp.ArpContext
	arpClient                        *arp.ArpClient
	xid                              [3]byte
	serverIp                         net.IP
	interfaceID                      string
	dhcpIn                           chan []byte
	preferredLifeTime, validLifeTime time.Time
	login                            string
}

func (self *dhcp6Context) resetLease() {
	if ipAddr := self.IpAddr.Load().(net.IP); !ipAddr.IsUnspecified() {
		arp.MapArpByIp.Reset(ipAddr)
		self.IpAddr.Store(net.IPv6unspecified)
	}

	self.serverIp = net.IPv6zero
}

func CreateClientv6(macAddr net.HardwareAddr, interfaceID, login string) (*Dhcpv6Client, chan []byte) {
	d := new(Dhcpv6Client)
	//TODO generate xid

	arpClient, arpContext := arp.ConstructArpClient(macAddr)
	d.ctx = &dhcp6Context{
		dhcpIn:      make(chan []byte, 100),
		arpClient:   arpClient,
		ArpContext:  arpContext,
		interfaceID: interfaceID,
		login:       login,
	}

	// At beginning,  the client send a SOLICIT
	d.currentState = solicitState{}

	return d, d.ctx.dhcpIn
}

type Dhcpv6Client struct {
	currentState iState
	ctx          *dhcp6Context
}

func (self *Dhcpv6Client) Run() {
	go func() {
		for {
			// Let's do the job forever...
			self.currentState = self.currentState.do(self.ctx)
		}
	}()
}