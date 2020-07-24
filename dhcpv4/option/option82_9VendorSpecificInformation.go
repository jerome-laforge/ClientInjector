package option

import (
	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/util"
)

type VendorSpecificInformation struct {
	EnterpriseNumber uint32
	Data             []byte
}

type Option82_9VendorSpecificInformation struct {
	rawOpt         RawOption
	vendorSpecInfo []VendorSpecificInformation
}

func (_ Option82_9VendorSpecificInformation) GetSubNum() byte {
	return byte(9)
}

func (this Option82_9VendorSpecificInformation) GetSubRawLen() int {
	return len(this.rawOpt.GetRawValue())
}

func (this Option82_9VendorSpecificInformation) GetSubLen() byte {
	return this.rawOpt.GetRawValue()[1]
}

func (this Option82_9VendorSpecificInformation) GetRawSubOption() RawOption {
	return this.rawOpt
}

func (this *Option82_9VendorSpecificInformation) Parse(opt RawOption) error {
	this.rawOpt = opt
	this.vendorSpecInfo = make([]VendorSpecificInformation, 0, 2)
	var length int
	value := opt.GetValue()
	for i := 0; i < len(value); {
		if i+5 > len(value) {
			return dherrors.Opt82_9IsMalformed
		}
		curVendorSpecInfo := new(VendorSpecificInformation)
		curVendorSpecInfo.EnterpriseNumber = util.Convert4byteToUint32(value[i : i+4])
		i += 4
		length = int(value[i])
		i++
		if i+length > len(value) {
			return dherrors.Opt82_9IsMalformed
		}
		curVendorSpecInfo.Data = value[i : i+length]
		i += length
		this.vendorSpecInfo = append(this.vendorSpecInfo, *curVendorSpecInfo)
	}
	return nil
}

func (this Option82_9VendorSpecificInformation) GetDataVendorSpecificInformation(enterpriseNumber uint32) ([]byte, bool) {
	for _, curVendorSpecInfo := range this.vendorSpecInfo {
		if curVendorSpecInfo.EnterpriseNumber == enterpriseNumber {
			return curVendorSpecInfo.Data, true
		}
	}
	return nil, false
}
