package option

type Option8QuotesServer struct {
	TypeListUint32
}

func (_ Option8QuotesServer) GetNum() byte {
	return byte(8)
}

func (this *Option8QuotesServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option8QuotesServer) Construct(quotesServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), quotesServer)
}

func (this Option8QuotesServer) GetQuotesServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option8QuotesServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
