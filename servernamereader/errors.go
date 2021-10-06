package servernamereader

import "errors"

var ErrNoSNI = errors.New("no SNI found in client hello")
var ErrSSL2 = errors.New("received SSL 2.0 Client Hello which can not support SNI")
var ErrNotTLSHandshake = errors.New("request did not begin with TLS handshake")
var ErrBadTLSVersion = errors.New("received SSL handshake which can not support SNI")
var ErrNotTLSHello = errors.New("not a client hello")
var ErrSSL3NoExtensions = errors.New("received SSL 3.0 handshake without extensions")
