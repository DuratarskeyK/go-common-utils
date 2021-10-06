package servernamereader

var protocolErrors = []error{
	ErrNoSNI,
	ErrSSL2,
	ErrNotTLSHandshake,
	ErrBadTLSVersion,
	ErrNotTLSHello,
	ErrSSL3NoExtensions,
}

func IsProtocolError(err error) bool {
	for _, v := range protocolErrors {
		if v == err {
			return true
		}
	}
	return false
}
