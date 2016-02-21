package main

import (
	"dhcpv4"
	"dhcpv4/option"
	"dhcpv4/util"
	"encoding/hex"
	"math/rand"
	"net"
	"time"
)

type DhcpClient struct {
	currentState iState
	ctx          dhcpContext
}

func (self *DhcpClient) run() {
	go func() {
		for {
			// Let's do the job forever...
			self.currentState = self.currentState.do()
		}
	}()
}

func (self DhcpClient) String() string {
	return "mac: " + self.ctx.macAddr.String() + " xid: 0x" + hex.EncodeToString(self.ctx.xid)
}

func CreateDhcpClient(macAddr net.HardwareAddr, giaddr uint32, login string) (*DhcpClient, error) {
	d := new(DhcpClient)

	arpClient, arpContext, err := ConstructArpClient(macAddr)
	if err != nil {
		return nil, err
	}

	xid := make([]byte, 4)
	util.ConvertUint32To4byte(rand.Uint32(), xid)
	d.ctx = dhcpContext{
		xid:        xid,
		dhcpIn:     make(chan []byte, 100),
		arpClient:  arpClient,
		ArpContext: arpContext,
		giaddr:     giaddr,
		login:      login,
	}

	// At beginning,  the client send a DISCOVER
	d.currentState = discoverState{
		dhcpContext: d.ctx,
	}

	if err != nil {
		return nil, err
	}

	return d, nil
}

type iState interface {
	do() iState
}

type dhcpContext struct {
	*ArpContext
	arpClient        ArpClient
	xid              []byte
	serverIp, giaddr uint32
	dhcpIn           chan []byte
	t1, t2, t0       time.Time
	state            iState
	login            string
}

func (self *dhcpContext) resetLease() {
	if ipAddr := self.ipAddr.Load().(uint32); ipAddr != 0 {
		dhcpContextByIp.ResetIp(ipAddr)
	}

	self.ipAddr.Store(uint32(0))
	self.serverIp = 0
}

func extractAllLeaseTime(dp *dhcpv4.DhcpPacket) (t0, t1, t2 time.Time) {
	now := time.Now()

	var durationT0, durationT1, durationT2 time.Duration
	opt51 := new(option.Option51IpAddressLeaseTime)
	if found, err := dp.GetOption(opt51); err == nil && found {
		durationT0 = time.Duration(opt51.GetLeaseTime()) * time.Second
	} else {
		durationT0 = 24 * time.Hour
	}

	if durationT0 < time.Minute {
		// fallback for t0
		durationT0 = 24 * time.Minute
	}

	opt58 := new(option.Option58RenewalTimeValue)
	if found, err := dp.GetOption(opt58); err == nil && found {
		durationT1 = time.Duration(opt58.GetRenewalTime()) * time.Second
	} else {
		durationT1 = durationT0 / 2
	}

	opt59 := new(option.Option59RebindingTimeValue)
	if found, err := dp.GetOption(opt59); err == nil && found {
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
