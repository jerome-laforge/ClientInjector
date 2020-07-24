package layer

import "github.com/google/gopacket"

type PayloadLayer struct {
	Contents []byte
	LT       gopacket.LayerType
}

func (self *PayloadLayer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.AppendBytes(len(self.Contents))
	if err != nil {
		return err
	}
	copy(bytes, self.Contents)
	return nil
}

// LayerType returns the type of the layer that is being serialized to the buffer
func (self *PayloadLayer) LayerType() gopacket.LayerType {
	return self.LT
}
