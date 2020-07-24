package option

type Option56Message struct {
	TypeString
}

func (_ Option56Message) GetNum() byte {
	return byte(56)
}

func (this *Option56Message) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option56Message) GetMessage() string {
	return this.TypeString.GetString()
}

func (this *Option56Message) Construct(message string) {
	this.TypeString.Construct(this.GetNum(), message)
}

func (this Option56Message) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
