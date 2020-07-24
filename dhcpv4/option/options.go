package option

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
)

type Option interface {
	GetType() byte
	GetLength() byte
	GetValue() []byte
}

type SpecificOption interface {
	GetNum() byte
	GetRawOption() RawOption
	Parse(RawOption) error
}

type RawOption struct {
	rawData []byte
}

func (opt RawOption) GetType() byte {
	return opt.rawData[0]
}

func (opt RawOption) GetLength() byte {
	if len(opt.rawData) < 2 {
		// for option like option #0 or option #255
		return 0
	}
	return opt.rawData[1]
}

func (opt RawOption) GetRawLen() int {
	return len(opt.rawData)
}

func (opt RawOption) GetValue() []byte {
	if opt.GetLength() == 0 {
		// for option like option #0 or option #255
		return nil
	}
	return opt.rawData[2:]
}

func (opt RawOption) ReplaceValue(value []byte) error {
	if len(value) != int(opt.GetLength()) {
		return errors.New(fmt.Sprintf("Option has not the same length [Actual : %d], [Expected : %d]", len(value), opt.GetLength()))
	}
	copy(opt.GetValue(), value)
	return nil
}

func (opt RawOption) GetIndexInto(includeRaw []byte) (int, error) {
	if len(includeRaw) == 0 {
		return -1, dherrors.BufferHas0Len
	}
	optionPointer := uintptr(unsafe.Pointer(&opt.rawData[0]))
	includeRawPointer := uintptr(unsafe.Pointer(&includeRaw[0]))

	idx := int(optionPointer - includeRawPointer)
	if idx < 0 || idx >= len(includeRaw) {
		return -2, dherrors.OptIsNotIntoBuf
	}

	return idx, nil
}

func (opt RawOption) GetRawValue() []byte {
	return opt.rawData
}

func (opt RawOption) CloneByte() []byte {
	cloneByte := make([]byte, len(opt.rawData))
	copy(cloneByte, opt.rawData)
	return cloneByte
}

func (opt *RawOption) Construct(typ, len byte) {
	if len == 0 {
		opt.rawData = make([]byte, 1)
	} else {
		opt.rawData = make([]byte, len+2)
		opt.rawData[1] = len
	}
	opt.rawData[0] = typ
}

func (opt RawOption) IsValid() bool {
	if opt.GetRawLen() == 0 {
		return false
	}

	if opt.GetType() == 0x00 || opt.GetType() == 0xff {
		return opt.GetLength() == 0
	}

	if opt.GetRawLen() < 2 {
		return false
	}

	return int(opt.GetLength()) == len(opt.GetValue())
}

func ParseOptions(dp []byte) ([]RawOption, error) {
	return ParseOptionsWithInitCap(dp, 6)
}

func ParseOptionsWithInitCap(dp []byte, initCap int) ([]RawOption, error) {
	options := make([]RawOption, 0, initCap)
	lenPacket := len(dp)
	var typ byte
	var len byte
	for idx := 0; idx < lenPacket; {
		typ = dp[idx]
		if typ != 0x00 && typ != 0xff {
			if idx+1 >= lenPacket {
				return nil, errors.New(fmt.Sprintf("Options malformed : not enough room for option %d's length, [expected: %d] [actual: %d]", typ, 1, lenPacket-idx))
			}
			len = dp[idx+1]
			if idx+int(len)+2 > lenPacket {
				return nil, errors.New(fmt.Sprintf("Options malformed : not enough room for option %d's data, [expected: %d] [actual: %d]", typ, int(len)+2, lenPacket-idx))
			}
			options = append(options, RawOption{
				rawData: dp[idx : idx+2+int(len) : idx+2+int(len)],
			})
			idx += int(len) + 2
		} else {
			options = append(options, RawOption{
				rawData: dp[idx : idx+1 : idx+1],
			})
			idx++

			if typ == 0xff {
				// After end option, subsequent octets should be filled with pad options. So ignore the rest of data if any
				break
			}
		}
	}
	return options, nil
}

func GetRawOption(options []RawOption, typeOpt byte) (v RawOption, b bool) {
	for _, v = range options {
		if v.GetType() == typeOpt {
			return v, true
		}
	}
	return v, false
}
