package option

import (
	"errors"
	"fmt"
)

type TypeString struct {
	rawOpt RawOption
}

func (tstring *TypeString) Parse(numOpt byte, rawOpt RawOption) error {
	tstring.rawOpt = rawOpt
	if tstring.rawOpt.GetLength() == 0 {
		return errors.New(fmt.Sprintf("Option %d has invalid null length", numOpt))
	}
	return nil
}

func (tstring *TypeString) GetRawOption() RawOption {
	return tstring.rawOpt
}

func (tstring TypeString) GetString() string {
	return string(tstring.rawOpt.GetValue()[:tstring.rawOpt.GetLength()])
}

func (tstring *TypeString) Construct(optNum byte, value string) {
	bArr := []byte(value)
	tstring.rawOpt.Construct(optNum, byte(len(bArr)))
	copy(tstring.rawOpt.GetValue(), bArr)
}

func (tstring *TypeString) GetType() byte {
	return tstring.rawOpt.GetType()
}

func (tstring *TypeString) GetLength() byte {
	return tstring.rawOpt.GetLength()
}
