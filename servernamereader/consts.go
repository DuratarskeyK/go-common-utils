package servernamereader

const tlsHeaderLen = 5
const tlsHandshakeContentType = 0x16
const tlsHandshakeTypeClientHello = 1

var TLSHandshakeFailureAlert = []byte{0x15, 0x03, 0x01, 0x00, 0x02, 0x02, 0x28}
