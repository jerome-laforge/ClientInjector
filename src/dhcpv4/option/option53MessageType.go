package option

import (
	"dhcpv4/dherrors"
	"errors"
	"fmt"
)

const LEN_OPT_53 = 1

type Option53MessageType struct {
	rawOpt RawOption
}

func (_ Option53MessageType) GetNum() byte {
	return byte(53)
}

func (this *Option53MessageType) Parse(rawOpt RawOption) error {
	if rawOpt.GetLength() != LEN_OPT_53 {
		return errors.New(fmt.Sprintf("Option 53 has invalid length [expected: 1] [actual: %d]", this.rawOpt.GetLength()))
	}
	this.rawOpt = rawOpt
	return nil
}

func (this Option53MessageType) GetLength() byte {
	return this.rawOpt.GetLength()
}

func (this Option53MessageType) GetValue() []byte {
	return this.rawOpt.GetValue()
}

func (this *Option53MessageType) Construct(msgType MessageType) error {
	if msgType == UNKNOWN {
		return dherrors.Opt53InvalidMessageType
	}
	this.rawOpt.Construct(this.GetNum(), LEN_OPT_53)
	this.rawOpt.GetValue()[0] = msgType.Value
	return nil
}

func (this Option53MessageType) GetRawOption() RawOption {
	return this.rawOpt
}

type MessageType struct {
	MessageType string
	Value       byte
}

func (this MessageType) String() string {
	return this.MessageType
}

var DHCPDISCOVER MessageType = MessageType{
	MessageType: "DHCPDISCOVER",
	Value:       1,
}

var DHCPOFFER MessageType = MessageType{
	MessageType: "DHCPOFFER",
	Value:       2,
}

var DHCPREQUEST MessageType = MessageType{
	MessageType: "DHCPREQUEST",
	Value:       3,
}

var DHCPDECLINE MessageType = MessageType{
	MessageType: "DHCPDECLINE",
	Value:       4,
}

var DHCPACK MessageType = MessageType{
	MessageType: "DHCPACK",
	Value:       5,
}

var DHCPNAK MessageType = MessageType{
	MessageType: "DHCPNAK",
	Value:       6,
}

var DHCPRELEASE MessageType = MessageType{
	MessageType: "DHCPRELEASE",
	Value:       7,
}

var DHCPINFORM MessageType = MessageType{
	MessageType: "DHCPINFORM",
	Value:       8,
}

var DHCPFORCERENEW MessageType = MessageType{
	MessageType: "DHCPFORCERENEW",
	Value:       9,
}

// RFC 4388
var DHCPLEASEQUERY MessageType = MessageType{
	MessageType: "DHCPLEASEQUERY",
	Value:       10,
}

// RFC 4388
var DHCPLEASEUNASSIGNED MessageType = MessageType{
	MessageType: "DHCPLEASEUNASSIGNED",
	Value:       11,
}

// RFC 4388
var DHCPLEASEUNKNOWN MessageType = MessageType{
	MessageType: "DHCPLEASEUNKNOWN",
	Value:       12,
}

// RFC 4388
var DHCPLEASEACTIVE MessageType = MessageType{
	MessageType: "DHCPLEASEACTIVE",
	Value:       13,
}

var UNKNOWN MessageType = MessageType{
	MessageType: "UNKNOWN",
	Value:       255,
}

func (this Option53MessageType) GetMessageType() MessageType {
	switch this.GetValue()[0] {
	case 1:
		return DHCPDISCOVER
	case 2:
		return DHCPOFFER
	case 3:
		return DHCPREQUEST
	case 4:
		return DHCPDECLINE
	case 5:
		return DHCPACK
	case 6:
		return DHCPNAK
	case 7:
		return DHCPRELEASE
	case 8:
		return DHCPINFORM
	case 9:
		return DHCPFORCERENEW
	case 10:
		return DHCPLEASEQUERY
	case 11:
		return DHCPLEASEUNASSIGNED
	case 12:
		return DHCPLEASEUNKNOWN
	case 13:
		return DHCPLEASEACTIVE
	default:
		return UNKNOWN
	}
}

func (this Option53MessageType) SetMessageType(messageType MessageType) error {
	return this.rawOpt.ReplaceValue([]byte{messageType.Value})
}
