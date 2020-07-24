package dherrors

type constant_dhcp_error string

func (this constant_dhcp_error) Error() string {
	return string(this)
}

// Errors
const (
	BufferHas0Len                      = constant_dhcp_error("buffer has 0 length")
	HardwareLenHastoEqual6For10mbEther = constant_dhcp_error("Hardware length has to be 6 for 10mb ethernet")
	MalformedOption                    = constant_dhcp_error("option is malformed")
	MissMagicCookie                    = constant_dhcp_error("Missing mandatory Magic Cookie")
	OptIsNotIntoBuf                    = constant_dhcp_error("option is not included into this buffer")
	OptNotFound                        = constant_dhcp_error("option not found")
	TooLongHardwareLen                 = constant_dhcp_error("Hardware length has to be lower than or equal to 16")
	Opt53InvalidMessageType            = constant_dhcp_error("Invalid message type")
	Opt53NotFound                      = constant_dhcp_error("Option 53 not found")
	Opt77InvalidLen                    = constant_dhcp_error("Option 77 has invalid length (< 2)")
	Opt82_9IsMalformed                 = constant_dhcp_error("option 82.9 is malformed")
	Opt90DataInvalidLen                = constant_dhcp_error("Option 90 : Data has invalid length")
	Opt120DataInvalidLen               = constant_dhcp_error("Data has invalid length")
	Opt120HasInvalidLen                = constant_dhcp_error("Option 120 has invalid length")
	Opt120UnsupportedEncoding          = constant_dhcp_error("Unsupported encoding")
	Opt121TooLong                      = constant_dhcp_error("Option 121 : too long")
	Opt255NotFound                     = constant_dhcp_error("Missing mandatory option END(255)")
)
