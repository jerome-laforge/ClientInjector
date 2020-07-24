package option

type Option60VendorClassIdentifier struct {
	TypeString
}

func (_ Option60VendorClassIdentifier) GetNum() byte {
	return byte(60)
}

func (this *Option60VendorClassIdentifier) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option60VendorClassIdentifier) GetVendorClassIdentifier() string {
	return this.TypeString.GetString()
}

func (this *Option60VendorClassIdentifier) Construct(vendorClassIdentrifier string) {
	this.TypeString.Construct(this.GetNum(), vendorClassIdentrifier)
}

func (this Option60VendorClassIdentifier) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
