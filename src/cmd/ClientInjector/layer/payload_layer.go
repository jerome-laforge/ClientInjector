package layer

import "github.com/google/gopacket"

type PayloadLayer struct {
	Contents []byte
}

func (self *PayloadLayer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.AppendBytes(len(self.Contents))
	if err != nil {
		return err
	}
	copy(bytes, self.Contents)
	return nil
}
