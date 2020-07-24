package option

import (
	"errors"
	"fmt"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/util"
)

const LEN_OPT_U_INT32 = 4

type TypeUint32 struct {
	rawOpt RawOption
}

func (tuint32 *TypeUint32) Parse(numOpt byte, rawOpt RawOption) error {
	tuint32.rawOpt = rawOpt
	if tuint32.GetLength() != LEN_OPT_U_INT32 {
		return errors.New(fmt.Sprintf("Option %d has invalid length [expected: 4] [actual: %d]", numOpt, tuint32.rawOpt.GetLength()))
	}
	return nil
}

func (tuint32 *TypeUint32) GetRawOption() RawOption {
	return tuint32.rawOpt
}

func (tuint32 TypeUint32) GetLength() byte {
	return tuint32.rawOpt.GetLength()
}

func (tuint32 TypeUint32) GetValue() []byte {
	return tuint32.rawOpt.GetValue()
}

func (tuint32 *TypeUint32) Construct(numOpt byte, number uint32) {
	tuint32.rawOpt.Construct(numOpt, LEN_OPT_U_INT32)
	util.ConvertUint32To4byte(number, tuint32.GetValue())
}

func (tuint32 TypeUint32) GetUint32() uint32 {
	return util.Convert4byteToUint32(tuint32.GetValue())
}
