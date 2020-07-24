package option

import (
	"errors"
	"fmt"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
)

type UserClassData struct {
	ucLen  byte
	ucData []byte
}

func (this UserClassData) GetUcLen() byte {
	return this.ucLen
}

func (this UserClassData) GetUcData() []byte {
	return this.ucData
}

func (this UserClassData) GetUcDataAsString() string {
	return string(this.ucData)
}

func (this *UserClassData) Construct(ucData []byte) {
	this.ucData = ucData
	this.ucLen = byte(len(ucData))
}

type Option77UserClass struct {
	rawOpt         RawOption
	UserClassDatas []UserClassData
}

func (_ Option77UserClass) GetNum() byte {
	return byte(77)
}

func (this *Option77UserClass) Parse(rawOpt RawOption) error {
	this.rawOpt = rawOpt
	if this.rawOpt.GetLength() < 2 {
		return dherrors.Opt77InvalidLen
	}
	if err := this.parseUserClassDataInstance(); err != nil {
		return err
	}
	return nil
}

func (this Option77UserClass) GetLength() byte {
	return this.rawOpt.GetLength()
}

func (this Option77UserClass) GetValue() []byte {
	return this.rawOpt.GetValue()
}

func (this *Option77UserClass) parseUserClassDataInstance() error {
	var idx byte = 0
	optData := this.GetValue()
	length := this.GetLength()
	var ucLen byte
	var ucData []byte
	for idx < length {
		ucLen = optData[idx]
		idx += 1
		if uint(ucLen)+uint(idx) > uint(length) {
			return errors.New(fmt.Sprintf("option 77 malformed : not enough room for UserClassData #%d [expected: %d] [actual: %d]", len(this.UserClassDatas), ucLen, length-idx))
		}
		ucData = optData[idx : idx+ucLen]
		this.UserClassDatas = append(this.UserClassDatas, UserClassData{
			ucLen:  ucLen,
			ucData: ucData,
		})
		idx += ucLen
	}
	return nil
}

func (this Option77UserClass) GetRawOption() RawOption {
	return this.rawOpt
}

func (this *Option77UserClass) Construct(ucs []UserClassData) {
	size := byte(0)
	for _, uc := range ucs {
		size += uc.GetUcLen() + 1
	}
	this.rawOpt.Construct(this.GetNum(), size)
	buf := this.rawOpt.GetValue()[:0]
	for _, uc := range ucs {
		buf = append(buf, uc.GetUcLen())
		buf = append(buf, uc.GetUcData()...)
	}
}
