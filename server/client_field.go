package server

import "bytes"

type (
	ClientField string
)

var (
	ClientFieldKey            ClientField = "key"
	ClientFieldHostname       ClientField = "hse"
	ClientFieldUserAgent      ClientField = "ura"
	ClientFieldConnectAddress ClientField = "cdr"
	ClientFieldRouteAddress   ClientField = "rdr"
	ClientFieldRemark         ClientField = "rem"
	ClientFieldOnline         ClientField = "ole"
	ClientFieldShell          ClientField = "she"
	ClientFieldStatus         ClientField = "sts"
	ClientFieldState          ClientField = "sta"
	ClientFieldCreatedAt      ClientField = "c"
	ClientFieldUpdatedAt      ClientField = "u"
	ClientFieldConnectAt      ClientField = "o"
	ClientFieldExpiredAt      ClientField = "e"

	AllClientFields = []ClientField{
		ClientFieldKey,
		ClientFieldHostname,
		ClientFieldConnectAddress,
		ClientFieldRouteAddress,
		ClientFieldRemark,
		ClientFieldOnline,
		ClientFieldShell,
		ClientFieldStatus,
		ClientFieldState,
		ClientFieldCreatedAt,
		ClientFieldUpdatedAt,
		ClientFieldConnectAt,
		ClientFieldExpiredAt,
	}
)

// client 的  dbkey
func clientFieldDBKey(id string, field ClientField) []byte {
	return bytes.Join([][]byte{[]byte("client"), []byte(field), []byte(id)}, []byte{'/'})
}
