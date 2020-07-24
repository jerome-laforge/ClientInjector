package option

import "github.com/jerome-laforge/ClientInjector/dhcpv4/util"

type Option119DomainSearch struct {
	rawOpt       RawOption
	domainSearch []string
}

func (_ Option119DomainSearch) GetNum() byte {
	return byte(119)
}

func (this *Option119DomainSearch) Parse(rawOpt RawOption) error {
	this.rawOpt = rawOpt
	var err error
	this.domainSearch, err = util.GetDnsName(this.rawOpt.GetValue())
	return err
}

func (this *Option119DomainSearch) Construct(fqdns []string) {
	data := util.SerializeDnsName(fqdns)
	this.rawOpt.Construct(this.GetNum(), byte(len(data)))
	copy(this.rawOpt.GetValue(), data)
}

func (this Option119DomainSearch) GetSearchString() []string {
	return this.domainSearch
}

func (this Option119DomainSearch) GetRawOption() RawOption {
	return this.rawOpt
}
