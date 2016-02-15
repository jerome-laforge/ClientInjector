# ClientInjector
ClientInjector is client that tries to emulate a real client (with MAC@ Spoofing, ARP).<br>
Currently, only DHCPv4 is supported (planed to add DHCPv6).<br>
Special thank to this wonderful project [gopacket] (https://github.com/google/gopacket)

## Build
Use [gb](http://getgb.io)

## Run
Testing only on Linux.
You have to be root for reading packet on wire (use c binding of libpcap via gopacket).
