package option

type Option82DhcpAgentOption struct {
	rawOpt    RawOption
	subOption []RawOption
}

type SubOption82 interface {
	GetSubNum() byte
	Parse(RawOption) error
	GetSubRawLen() int
	GetSubLen() byte
	GetRawSubOption() RawOption
}

func (_ Option82DhcpAgentOption) GetNum() byte {
	return byte(82)
}

func (this *Option82DhcpAgentOption) Parse(rawOpt RawOption) error {
	this.rawOpt = rawOpt
	var err error
	this.subOption, err = ParseOptionsWithInitCap(this.rawOpt.GetValue(), 2)
	return err
}

func (this Option82DhcpAgentOption) GetSubOptions(specSubOpt SubOption82) (bool, error) {
	genOpt, found := GetRawOption(this.subOption, specSubOpt.GetSubNum())
	if found {
		err := specSubOpt.Parse(genOpt)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (this Option82DhcpAgentOption) GetLength() byte {
	return this.rawOpt.GetLength()
}

func (this Option82DhcpAgentOption) GetValue() []byte {
	return this.rawOpt.GetValue()
}

func (this Option82DhcpAgentOption) GetRawOption() RawOption {
	return this.rawOpt
}

func (this *Option82DhcpAgentOption) Construct(subOpts []SubOption82) {
	size := 0
	for _, subOpt := range subOpts {
		size += subOpt.GetSubRawLen()
	}

	this.rawOpt.Construct(this.GetNum(), byte(size))
	buf := this.rawOpt.GetValue()[:0]
	for _, subOpt := range subOpts {
		buf = append(buf, subOpt.GetRawSubOption().GetRawValue()...)
	}
}
