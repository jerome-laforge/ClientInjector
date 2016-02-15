package option

type Option49XWindowSystemDisplayManager struct {
	TypeListUint32
}

func (_ Option49XWindowSystemDisplayManager) GetNum() byte {
	return byte(49)
}

func (this *Option49XWindowSystemDisplayManager) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option49XWindowSystemDisplayManager) Construct(xWindowSystemDisplayManager []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), xWindowSystemDisplayManager)
}

func (this Option49XWindowSystemDisplayManager) GetXWindowSystemDisplayManager() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option49XWindowSystemDisplayManager) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
