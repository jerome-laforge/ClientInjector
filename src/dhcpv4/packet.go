package dhcpv4

import (
	"bytes"
	"errors"
	"fmt"

	"dhcpv4/dherrors"
	"dhcpv4/option"
	"dhcpv4/util"
)

type DhcpPacket struct {
	Raw     []byte
	Options []option.RawOption
}

func (dp *DhcpPacket) ConstructWithPreAllocatedBuffer(buffer []byte, msgType option.MessageType) (err error) {
	if cap(buffer) < 1500 {
		buffer = make([]byte, 241, 1500)
	} else if len(buffer) != 241 {
		buffer = buffer[:241]
	}
	dp.Raw = buffer
	dp.SetHardwareType(byte(1))
	dp.SetHardwareLength(byte(6))
	dp.SetMagicCookie()
	err = dp.ResetAllOptions()
	if err != nil {
		return err
	}
	opt53 := new(option.Option53MessageType)
	opt53.Construct(msgType)
	dp.AddOption(opt53)

	switch msgType {
	case option.DHCPDISCOVER:
		fallthrough
	case option.DHCPREQUEST:
		fallthrough
	case option.DHCPDECLINE:
		fallthrough
	case option.DHCPRELEASE:
		dp.SetOpCode(OP_BOOTREQUEST)
		// break
	case option.DHCPOFFER:
		fallthrough
	case option.DHCPACK:
		fallthrough
	case option.DHCPNAK:
		dp.SetOpCode(OP_BOOTREPLY)
		// break
	}
	return nil
}

func (dp *DhcpPacket) Construct(msgType option.MessageType) {
	dp.ConstructWithPreAllocatedBuffer(make([]byte, 241, 1500), msgType)
}

func Parse(raw []byte) (dp *DhcpPacket, err error) {
	dp = new(DhcpPacket)
	dp.Raw = raw
	err = dp.headerSanityCheck()
	if err != nil {
		dp = nil
		return
	}

	dp.Options, err = option.ParseOptions(raw[240:]) //240th byte is where options are
	if err != nil {
		dp = nil
		return
	}
	if err = dp.optionSanityCheck(); err != nil {
		dp = nil
		return
	}
	return
}

var magic_cookie = []byte{0x63, 0x82, 0x53, 0x63}

type OpCode struct {
	OpCode string
	Value  byte
}

var OP_BOOTREQUEST OpCode = OpCode{
	OpCode: "BOOTREQUEST",
	Value:  1,
}

var OP_BOOTREPLY OpCode = OpCode{
	OpCode: "BOOTREPLY",
	Value:  2,
}

var UNKNOWN OpCode = OpCode{
	OpCode: "UNKNOWN",
	Value:  255,
}

func (dp DhcpPacket) GetOpCode() OpCode {
	switch dp.Raw[0] {
	case 1:
		return OP_BOOTREQUEST
	case 2:
		return OP_BOOTREPLY
	default:
		return UNKNOWN
	}
}

func (dp DhcpPacket) SetOpCode(opCode OpCode) {
	dp.Raw[0] = opCode.Value
}

func (dp DhcpPacket) GetHardwareType() byte {
	return dp.Raw[1]
}

func (dp DhcpPacket) SetHardwareType(hardwareType byte) {
	dp.Raw[1] = hardwareType
}

func (dp DhcpPacket) GetHardwareLength() byte {
	return dp.Raw[2]
}

func (dp DhcpPacket) SetHardwareLength(hardwareLength byte) {
	dp.Raw[2] = hardwareLength
}

func (dp DhcpPacket) GetHops() byte {
	return dp.Raw[3]
}

func (dp DhcpPacket) SetHops(hops byte) {
	dp.Raw[3] = hops
}

func (dp DhcpPacket) GetXid() uint32 {
	return util.Convert4byteToUint32(dp.Raw[4:8])
}

func (dp DhcpPacket) GetXidRaw() []byte {
	return dp.Raw[4:8:8]
}

func (dp DhcpPacket) SetXid(xid uint32) {
	util.ConvertUint32To4byte(xid, dp.Raw[4:8])
}

func (dp DhcpPacket) GetClientIp() uint32 {
	return util.Convert4byteToUint32(dp.Raw[12:16])
}

func (dp DhcpPacket) SetClientIp(ip uint32) {
	util.ConvertUint32To4byte(ip, dp.Raw[12:16])
}

func (dp DhcpPacket) GetYourIp() uint32 {
	return util.Convert4byteToUint32(dp.Raw[16:20])
}

func (dp DhcpPacket) SetYourIp(ip uint32) {
	util.ConvertUint32To4byte(ip, dp.Raw[16:20])
}

func (dp DhcpPacket) GetNextServerIp() uint32 {
	return util.Convert4byteToUint32(dp.Raw[20:24])
}

func (dp DhcpPacket) SetNextServerIp(ip uint32) {
	util.ConvertUint32To4byte(ip, dp.Raw[20:24])
}

func (dp DhcpPacket) GetGiAddr() uint32 {
	return util.Convert4byteToUint32(dp.Raw[24:28])
}

func (dp DhcpPacket) SetGiAddr(ip uint32) {
	util.ConvertUint32To4byte(ip, dp.Raw[24:28])
}

func (dp DhcpPacket) GetMacAddr() []byte {
	hi := 28 + int(dp.GetHardwareLength())
	return dp.Raw[28:hi:hi]
}

func (dp DhcpPacket) SetMacAddr(macAddr []byte) {
	copy(dp.Raw[28:28+int(dp.GetHardwareLength())], macAddr)
}

func (dp DhcpPacket) ContainsMagicCookie() bool {
	return bytes.Equal(magic_cookie, dp.Raw[236:240])
}

func (dp DhcpPacket) SetMagicCookie() {
	copy(dp.Raw[236:240], magic_cookie)
}

func (dp DhcpPacket) optionSanityCheck() error {
	if dp.Options[len(dp.Options)-1].GetType() != 0xff {
		return dherrors.Opt255NotFound
	}
	if !dp.ContainOption(new(option.Option53MessageType)) {
		return dherrors.Opt53NotFound
	}
	return nil
}

func (dp DhcpPacket) headerSanityCheck() error {
	//For DHCP : BOOTP head's length equals 240 bytes
	//Plus 3 bytes for option 53 (as 53 is mandatory option)
	if len(dp.Raw) < 243 {
		return errors.New(fmt.Sprintf("Length has to be equal or greater than 243 bytes [actual:%d]", len(dp.Raw)))
	}
	if !dp.ContainsMagicCookie() {
		return dherrors.MissMagicCookie
	}
	if dp.GetHardwareLength() > 16 {
		return dherrors.TooLongHardwareLen
	}
	if dp.GetHardwareType() == byte(1) && dp.GetHardwareLength() != byte(6) {
		return dherrors.HardwareLenHastoEqual6For10mbEther
	}
	return nil
}

func (dp DhcpPacket) GetOption(specOpt option.SpecificOption) (bool, error) {
	genOpt, found := option.GetRawOption(dp.Options, specOpt.GetNum())
	if found {
		err := specOpt.Parse(genOpt)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (dp DhcpPacket) GetRawOption(typeOpt byte) (option.RawOption, bool) {
	return option.GetRawOption(dp.Options, typeOpt)
}

func (dp *DhcpPacket) ResetAllOptions() error {
	var err error
	dp.Raw[240] = 0xff
	dp.Raw = dp.Raw[:241]                               // remove padding and hide it.
	dp.Options, err = option.ParseOptions(dp.Raw[240:]) //240th byte is where options are
	return err
}

func (dp *DhcpPacket) RemoveAllOptionsExcept(numOptList []byte) (err error) {
ignoreThisOption:
	for i := len(dp.Options) - 2; 0 <= i; i-- {
		for _, numOpt := range numOptList {
			if dp.Options[i].GetType() == numOpt {
				continue ignoreThisOption
			}
		}

		rawOpt := dp.Options[i]
		// test whether this option is included into raw packet's data
		idx, err := rawOpt.GetIndexInto(dp.Raw)
		if err != nil {
			// this option is not included ???
			continue ignoreThisOption
		}

		// remove data from rawData
		dp.Raw = append(dp.Raw[:idx], dp.Raw[idx+int(rawOpt.GetRawLen()):]...)
	}
	dp.Options, err = option.ParseOptions(dp.Raw[240:]) //240th byte is where options are
	return
}

func (dp *DhcpPacket) AddOption(optionToAdd option.SpecificOption) error {
	return dp.AddRawOption(optionToAdd.GetRawOption())
}

func (dp *DhcpPacket) AddRawOption(optionToAdd option.RawOption) error {
	if !optionToAdd.IsValid() {
		return dherrors.MalformedOption
	}

	idx, err := dp.Options[len(dp.Options)-1].GetIndexInto(dp.Raw)
	if err != nil {
		// this option is not included ???
		return err
	}

	// add data from rawData
	dp.Raw = append(dp.Raw[:idx], append(optionToAdd.GetRawValue(), dp.Raw[idx:]...)...)
	dp.Options, err = option.ParseOptionsWithInitCap(dp.Raw[240:], len(dp.Options)+1) //240th byte is where options are
	return err
}

func (dp *DhcpPacket) AddOptionBefore(optionToAdd, beforeOption option.SpecificOption) (err error) {
	if !optionToAdd.GetRawOption().IsValid() {
		return dherrors.MalformedOption
	}

	done := false
	for i := len(dp.Options) - 2; 0 <= i; i-- {
		if dp.Options[i].GetType() != beforeOption.GetNum() {
			continue
		}
		if 0 < i && dp.Options[i-1].GetType() == beforeOption.GetNum() {
			// we don't insert here as this option use RFC 3396 : Encoding Long Options in the Dynamic Host Configuration Protocol (DHCPv4)
			continue
		}
		done = true
		rawOpt := dp.Options[i]
		// test whether this option is included into raw packet's data
		idx, err := rawOpt.GetIndexInto(dp.Raw)
		if err != nil {
			// this option is not included ???
			continue
		}

		// insert data before option into rawData
		dp.Raw = append(dp.Raw[:idx], append(optionToAdd.GetRawOption().GetRawValue(), dp.Raw[idx:]...)...)
		break
	}
	if done {
		dp.Options, err = option.ParseOptions(dp.Raw[240:]) //240th byte is where options are
	} else {
		err = dherrors.OptNotFound
	}
	return
}

func (dp *DhcpPacket) RemoveOption(numOpt byte) (err error) {
	for i := len(dp.Options) - 2; 0 <= i; i-- {
		if dp.Options[i].GetType() == numOpt {
			rawOpt := dp.Options[i]
			// test whether this option is included into raw packet's data
			idx, err := rawOpt.GetIndexInto(dp.Raw)
			if err != nil {
				// this option is not included ???
				continue
			}

			// remove data from rawData
			dp.Raw = append(dp.Raw[:idx], dp.Raw[idx+int(rawOpt.GetRawLen()):]...)
		}
	}
	dp.Options, err = option.ParseOptions(dp.Raw[240:]) //240th byte is where options are
	return
}

func (dp *DhcpPacket) ContainOption(optionToCheck option.SpecificOption) bool {
	for _, cOption := range dp.Options {
		if cOption.GetType() == optionToCheck.GetNum() {
			return true
		}
	}
	return false
}

func (dp *DhcpPacket) ReplaceOption(optionToReplace option.SpecificOption) (err error) {
	done := false
	for i := len(dp.Options) - 2; 0 <= i; i-- {
		if dp.Options[i].GetType() == optionToReplace.GetNum() {
			rawOpt := dp.Options[i]
			// test whether this option is included into raw packet's data
			idx, err := rawOpt.GetIndexInto(dp.Raw)
			if err != nil {
				// this option is not included ???
				continue
			}

			if done {
				// remove old data option from rawData.
				// It is used when several time same option in same packet
				// For example RFC 3396 : Encoding Long Options in the Dynamic Host Configuration Protocol (DHCPv4)
				dp.Raw = append(dp.Raw[:idx], dp.Raw[idx+int(rawOpt.GetRawLen()):]...)
			} else {
				// insert option into rawData
				dp.Raw = append(dp.Raw[:idx], append(optionToReplace.GetRawOption().GetRawValue(), dp.Raw[idx+int(rawOpt.GetRawLen()):]...)...)
				done = true
			}
		}
	}
	if done {
		dp.Options, err = option.ParseOptions(dp.Raw[240:]) //240th byte is where options are
	} else {
		err = dherrors.OptNotFound
	}
	return
}

func (dp DhcpPacket) GetNumberOfOption() int {
	return len(dp.Options)
}

func (dp DhcpPacket) GetOptionAt(idx int) (option.RawOption, error) {
	if idx < 0 || len(dp.Options)-1 < idx {
		return option.RawOption{}, errors.New(fmt.Sprintf("Out of bound option for index %d. Valid range [0, %d]", idx, len(dp.Options)-1))
	}
	return dp.Options[idx], nil
}

func (dp DhcpPacket) GetPadding() []byte {
	if idx, err := dp.Options[len(dp.Options)-1].GetIndexInto(dp.Raw); err == nil {
		return dp.Raw[idx+1:]
	}
	return []byte{}
}

func (dp DhcpPacket) GetTypeMessage() (option.MessageType, error) {
	opt53 := new(option.Option53MessageType)
	found, err := dp.GetOption(opt53)
	if err != nil {
		return option.UNKNOWN, err
	}
	if !found {
		return option.UNKNOWN, dherrors.Opt53NotFound
	}
	return opt53.GetMessageType(), nil
}
