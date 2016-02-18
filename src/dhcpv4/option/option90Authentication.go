package option

import (
	"dhcpv4/dherrors"
	"dhcpv4/util"
)

const HEADER_LEN_OPT_90 = 12

type Option90Authentificiation struct {
	rawOpt RawOption
}

func (this *Option90Authentificiation) Construct(data []byte) error {
	if len(data) < HEADER_LEN_OPT_90 {
		return dherrors.Opt90DataInvalidLen
	}
	this.rawOpt.Construct(this.GetNum(), byte(len(data)))

	copy(this.rawOpt.GetValue(), data)
	return nil
}

func (_ Option90Authentificiation) GetNum() byte {
	return byte(90)
}

func (this *Option90Authentificiation) Parse(rawOpt RawOption) error {
	if rawOpt.GetLength() <= HEADER_LEN_OPT_90 {
		return dherrors.MalformedOption
	}
	this.rawOpt = rawOpt
	return nil
}

func (this Option90Authentificiation) GetRawOption() RawOption {
	return this.rawOpt
}

func (this Option90Authentificiation) GetProtocol() byte {
	return this.rawOpt.GetValue()[0]
}

func (this Option90Authentificiation) GetAlgorithm() byte {
	return this.rawOpt.GetValue()[1]
}

func (this Option90Authentificiation) GetRdm() byte {
	return this.rawOpt.GetValue()[2]
}

func (this Option90Authentificiation) GetReplayDetection() uint64 {
	return util.Convert8byteToUint64(this.rawOpt.GetValue()[3:11])
}

func (this Option90Authentificiation) GetAuthenticationInformation() []byte {
	return this.rawOpt.GetValue()[11:]
}
