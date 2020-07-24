package option

type Option82_2RemoteId struct {
	rawOpt RawOption
}

func (_ Option82_2RemoteId) GetSubNum() byte {
	return byte(2)
}

func (this Option82_2RemoteId) GetSubRawLen() int {
	return len(this.rawOpt.GetRawValue())
}

func (this Option82_2RemoteId) GetSubLen() byte {
	return this.rawOpt.GetRawValue()[1]
}

func (this *Option82_2RemoteId) Parse(opt RawOption) error {
	this.rawOpt = opt
	return nil
}

func (this Option82_2RemoteId) GetRemoteId() []byte {
	return this.rawOpt.GetValue()
}

func (this *Option82_2RemoteId) Construct(value []byte) {
	this.rawOpt.Construct(this.GetSubNum(), byte(len(value)))
	this.rawOpt.ReplaceValue(value)
}

func (this Option82_2RemoteId) GetRawSubOption() RawOption {
	return this.rawOpt
}
