package llrp

const (
	// LLRPStatus StatusCode
	// message scope
	M_Success        = 0
	M_ParameterError = 100 + iota
	M_FieldError
	M_UnexpectedParameter
	M_MissingParameter
	M_DuplicateParameter
	M_OverflowParameter
	M_OverflowField
	M_UnknownParameter
	M_UnknownField
	M_UnsupportedMessage
	M_UnsupportedVersion
	M_UnsupportedParameter
	M_UnexpectedMessage
)
const (
	// parameter
	P_ParameterError_Scope = 200 + iota
	P_FieldError_Scope
	P_UnexpectedParameter
	P_MissingParameter
	P_DuplicateParamete
	P_OverflowParameter
	P_OverflowField
	P_UnknownParameter
	P_UnknownField
	P_UnsupportedParameter
)
const (
	// field
	A_Invalid = 300 + iota
	A_OutOfRange
)
const (
	// reader
	R_DeviceError = 401
)
