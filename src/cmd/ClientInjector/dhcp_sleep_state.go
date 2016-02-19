package main

import (
	"dhcpv4/util"
	"log"
	"net"
	"time"
)

type sleepState struct {
	dhcpContext
}

func (self sleepState) do() iState {
	var (
		macAddr = self.macAddr.Load().(net.HardwareAddr)
		ipAddr  = self.ipAddr.Load().(uint32)
	)

	log.Println(macAddr, "ip", util.ConvertUint32ToIpAddr(ipAddr), "sleep until t1", self.t1)
	time.Sleep(self.t1.Sub(time.Now()))

	return requestRenewState{
		dhcpContext: self.dhcpContext,
	}
}

type timeoutRenewState struct {
	dhcpContext
}

func (self timeoutRenewState) do() iState {
	var (
		macAddr   = self.macAddr.Load().(net.HardwareAddr)
		ipAddr    = self.ipAddr.Load().(uint32)
		now       = time.Now()
		timeout   = self.t2.Sub(now) / 2
		nextState iState
	)

	if timeout < time.Minute {
		timeout = self.t2.Sub(now)
		nextState = &requestRebindState{
			dhcpContext: self.dhcpContext,
		}
	} else {
		nextState = &requestRenewState{
			dhcpContext: self.dhcpContext,
		}
	}

	log.Println(macAddr, "ip", util.ConvertUint32ToIpAddr(ipAddr), "sleep until ", now.Add(timeout))
	time.Sleep(timeout)
	return nextState
}

type timeoutRebindState struct {
	dhcpContext
}

func (self timeoutRebindState) do() iState {
	var (
		macAddr   = self.macAddr.Load().(net.HardwareAddr)
		ipAddr    = self.ipAddr.Load().(uint32)
		now       = time.Now()
		timeout   = self.t0.Sub(now) / 2
		nextState iState
	)

	if timeout < time.Minute {
		timeout = self.t0.Sub(now)
		nextState = &discoverState{
			dhcpContext: self.dhcpContext,
		}
	} else {
		nextState = &requestRebindState{
			dhcpContext: self.dhcpContext,
		}
	}

	log.Println(macAddr, "ip", util.ConvertUint32ToIpAddr(ipAddr), "sleep until ", now.Add(timeout))
	time.Sleep(timeout)
	return nextState
}
