package option

import (
	"dhcpv4/util"
	"errors"
	"fmt"
)

type TypeListUint32 struct {
	rawOpt RawOption
}

func (typ *TypeListUint32) Parse(numOpt byte, rawOpt RawOption) error {
	typ.rawOpt = rawOpt
	if typ.rawOpt.GetLength() < 4 || typ.rawOpt.GetLength()&3 != 0 {
		return errors.New(fmt.Sprintf("Option %d has invalid length", numOpt))
	}
	return nil
}

func (typ *TypeListUint32) GetRawOption() RawOption {
	return typ.rawOpt
}

func (typ TypeListUint32) GetListUint32() []uint32 {
	raw := typ.rawOpt.GetValue()
	listUint32 := make([]uint32, len(raw)/int(LEN_OPT_U_INT32))
	for i := range listUint32 {
		listUint32[i] = util.Convert4byteToUint32(raw[i*int(LEN_OPT_U_INT32) : (i+1)*int(LEN_OPT_U_INT32)])
	}
	return listUint32
}

func (typ *TypeListUint32) Construct(optNum byte, value []uint32) error {
	if len(value) == 0 {
		return errors.New(fmt.Sprintf("Option %d has invalid null length", optNum))
	}
	typ.rawOpt.Construct(optNum, byte(len(value)*int(LEN_OPT_U_INT32)))
	optData := typ.rawOpt.GetValue()
	for i, v := range value {
		util.ConvertUint32To4byte(v, optData[i*int(LEN_OPT_U_INT32):(i+1)*int(LEN_OPT_U_INT32)])
	}
	return nil
}

func (typ *TypeListUint32) GetType() byte {
	return typ.rawOpt.GetType()
}

func (typ *TypeListUint32) GetLength() byte {
	return typ.rawOpt.GetLength()
}

func (typ *TypeListUint32) GetValue() []byte {
	return typ.rawOpt.GetValue()
}
