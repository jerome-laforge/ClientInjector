package main

import "github.com/google/gopacket"

type PayloadLayer struct {
	contents []byte
}

func (self *PayloadLayer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.AppendBytes(len(self.contents))
	if err != nil {
		return err
	}
	copy(bytes, self.contents)
	return nil
}
