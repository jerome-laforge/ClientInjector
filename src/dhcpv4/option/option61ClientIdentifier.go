package option

import "dhcpv4/dherrors"

const HEADER_LEN_OPT_61 = 1

type Option61ClientIdentifier struct {
	rawOpt RawOption
}

func (_ Option61ClientIdentifier) GetNum() byte {
	return byte(61)
}

func (this *Option61ClientIdentifier) Parse(rawOpt RawOption) error {
	if rawOpt.GetLength() <= HEADER_LEN_OPT_61 {
		return dherrors.MalformedOption
	}
	this.rawOpt = rawOpt
	return nil
}

func (this *Option61ClientIdentifier) Construct(typ byte, clientIdentifier []byte) {
	this.rawOpt.Construct(this.GetNum(), HEADER_LEN_OPT_61+byte(len(clientIdentifier)))
	this.rawOpt.GetValue()[0] = typ
	copy(this.GetClientIdentifier(), clientIdentifier)
}

func (this Option61ClientIdentifier) GetRawOption() RawOption {
	return this.rawOpt
}

func (this Option61ClientIdentifier) GetType() byte {
	return this.rawOpt.GetValue()[0]
}

func (this Option61ClientIdentifier) GetClientIdentifier() []byte {
	return this.rawOpt.GetValue()[1:]
}
