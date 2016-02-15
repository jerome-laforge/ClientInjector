package option

type Option82_1CircuitId struct {
	rawOpt RawOption
}

func (_ Option82_1CircuitId) GetSubNum() byte {
	return byte(1)
}

func (this Option82_1CircuitId) GetSubRawLen() int {
	return len(this.rawOpt.GetRawValue())
}

func (this Option82_1CircuitId) GetSubLen() byte {
	return this.rawOpt.GetRawValue()[1]
}

func (this *Option82_1CircuitId) Parse(opt RawOption) error {
	this.rawOpt = opt
	return nil
}

func (this Option82_1CircuitId) GetCircuitId() []byte {
	return this.rawOpt.GetValue()
}

func (this *Option82_1CircuitId) Construct(value []byte) {
	this.rawOpt.Construct(this.GetSubNum(), byte(len(value)))
	this.rawOpt.ReplaceValue(value)
}

func (this Option82_1CircuitId) GetRawSubOption() RawOption {
	return this.rawOpt
}
