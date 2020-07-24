package dhcp4

import (
	"encoding/hex"
	"math/rand"
	"net"
	"time"

	"github.com/jerome-laforge/ClientInjector/cmd/ClientInjector/arp"
	"github.com/jerome-laforge/ClientInjector/dhcpv4"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/option"
)

var (
	DhcRelay bool
	Option90 bool
)

type DhcpClient struct {
	currentState iState
	ctx          *dhcpContext
}

func (self *DhcpClient) Run() {
	go func() {
		for {
			// Let's do the job forever...
			self.currentState = self.currentState.do(self.ctx)
		}
	}()
}

func (self DhcpClient) String() string {
	return "mac: " + self.ctx.MacAddr.String() + " xid: 0x" + hex.EncodeToString(self.ctx.xid[:])
}

func CreateClient(macAddr net.HardwareAddr, giaddr uint32, login string) (*DhcpClient, chan []byte) {
	d := new(DhcpClient)

	arpClient, arpContext := arp.ConstructArpClient(macAddr)

	d.ctx = &dhcpContext{
		dhcpIn:     make(chan []byte, 100),
		arpClient:  arpClient,
		ArpContext: arpContext,
		giaddr:     giaddr,
		login:      login,
	}

	rand.Read(d.ctx.xid[:])

	// At beginning,  the client send a DISCOVER
	d.currentState = discoverState{}

	return d, d.ctx.dhcpIn
}

type iState interface {
	do(*dhcpContext) iState
}

type dhcpContext struct {
	*arp.ArpContext
	arpClient        *arp.ArpClient
	xid              [4]byte
	serverIp, giaddr uint32
	dhcpIn           chan []byte
	t0, t1, t2       time.Time
	login            string
}

func (self *dhcpContext) resetLease() {
	self.serverIp = 0

	if ipAddr := self.IpAddr.Load().(net.IP); !ipAddr.IsUnspecified() {
		rand.Read(self.xid[:])
		arp.MapArpByIp.Reset(ipAddr)
		self.IpAddr.Store(net.IPv4zero)
	}
}

func extractAllLeaseTime(dp *dhcpv4.DhcpPacket) (t0, t1, t2 time.Time) {
	now := time.Now()

	var durationT0, durationT1, durationT2 time.Duration
	opt51 := new(option.Option51IpAddressLeaseTime)
	if found, _ := dp.GetOption(opt51); found {
		durationT0 = time.Duration(opt51.GetLeaseTime()) * time.Second
	} else {
		durationT0 = 24 * time.Hour
	}

	if durationT0 < time.Minute {
		// fallback for t0
		durationT0 = 24 * time.Minute
	}

	opt58 := new(option.Option58RenewalTimeValue)
	if found, _ := dp.GetOption(opt58); found {
		durationT1 = time.Duration(opt58.GetRenewalTime()) * time.Second
	} else {
		durationT1 = durationT0 / 2
	}

	opt59 := new(option.Option59RebindingTimeValue)
	if found, _ := dp.GetOption(opt59); found {
		durationT2 = time.Duration(opt59.GetRebindingTime()) * time.Second
	} else {
		durationT2 = durationT0 * 7 / 8
	}

	if !(durationT1 < durationT2 && durationT2 < durationT0) {
		// fallback for t1 & t2
		durationT1 = durationT0 / 2
		durationT2 = durationT0 * 7 / 8
	}

	t0 = now.Add(durationT0)
	t1 = now.Add(durationT1)
	t2 = now.Add(durationT2)
	return
}

func generateOption82(macAddr net.HardwareAddr) *option.Option82DhcpAgentOption {
	opt82_1 := new(option.Option82_1CircuitId)
	opt82_1.Construct([]byte(hex.EncodeToString(macAddr)))

	opt82_2 := new(option.Option82_2RemoteId)
	opt82_2.Construct([]byte(hex.EncodeToString(macAddr)))

	opt82 := new(option.Option82DhcpAgentOption)
	opt82.Construct([]option.SubOption82{
		opt82_1,
		opt82_2,
	})

	return opt82
}

func generateOption90(login string) *option.Option90Authentificiation {
	opt90 := new(option.Option90Authentificiation)
	buf := make([]byte, option.HEADER_LEN_OPT_90+len(login))
	copy(buf[option.HEADER_LEN_OPT_90:], []byte(login))
	opt90.Construct(buf)

	return opt90
}
