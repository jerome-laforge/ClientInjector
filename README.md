# ClientInjector
ClientInjector is client that tries to emulate a real client (with MAC@ Spoofing, ARP).<br>
Currently, only DHCPv4 is supported (planed to add DHCPv6).<br>
Special thank to this wonderful project [gopacket](https://github.com/google/gopacket)<br>

## Build
Use [gb](http://getgb.io)<br>
Pcap dev headers might be necessary<br>
```
apt-get install libpcap-dev
```

## Run
Testing only on Linux.<br>
You have to be root for reading packet on wire (use c binding of libpcap via gopacket).<br>
