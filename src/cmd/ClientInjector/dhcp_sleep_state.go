package main

import (
	"dhcpv4/util"
	"log"
	"time"
)

type sleepState struct {
	dhcpContext
}

func (self sleepState) do() iState {
	log.Println(self.macAddr, "ip", util.ConvertUint32ToIpAddr(self.ipAddr), "sleep until t1", self.t1)
	time.Sleep(self.t1.Sub(time.Now()))

	return requestSelectState{
		dhcpContext: self.dhcpContext,
	}
}
