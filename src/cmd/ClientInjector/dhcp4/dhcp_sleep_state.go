package dhcp4

import (
	"log"
	"net"
	"time"
)

type sleepState struct{}

func (_ sleepState) do(ctx *dhcpContext) iState {
	ipAddr := ctx.IpAddr.Load().(net.IP)

	log.Println(ctx.MacAddr, "ip", ipAddr, "sleep until t1", ctx.t1)
	time.Sleep(ctx.t1.Sub(time.Now()))

	return requestRenewState{}
}

type timeoutRenewState struct{}

func (_ timeoutRenewState) do(ctx *dhcpContext) iState {
	var (
		ipAddr    = ctx.IpAddr.Load().(net.IP)
		now       = time.Now()
		timeout   = ctx.t2.Sub(now) / 2
		nextState iState
	)

	if timeout < time.Minute {
		timeout = ctx.t2.Sub(now)
		nextState = requestRebindState{}
	} else {
		nextState = requestRenewState{}
	}

	log.Println(ctx.MacAddr, "ip", ipAddr, "sleep until ", now.Add(timeout))
	time.Sleep(timeout)
	return nextState
}

type timeoutRebindState struct{}

func (_ timeoutRebindState) do(ctx *dhcpContext) iState {
	var (
		ipAddr    = ctx.IpAddr.Load().(net.IP)
		now       = time.Now()
		timeout   = ctx.t0.Sub(now) / 2
		nextState iState
	)

	if timeout < time.Minute {
		// lease will be expired because DHCP Client didn't receive ACK for all its REQUEST.
		// DHCP Client will sent DISCOVER at the end of its lease.
		timeout = ctx.t0.Sub(now)
		nextState = discoverState{}
	} else {
		nextState = requestRebindState{}
	}

	log.Println(ctx.MacAddr, "ip", ipAddr, "sleep until ", now.Add(timeout))
	time.Sleep(timeout)
	return nextState
}
