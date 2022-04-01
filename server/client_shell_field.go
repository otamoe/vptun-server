package server

import "bytes"

type (
	ClientShellField string
)

var (
	ClientShellFieldInput     ClientShellField = "int"
	ClientShellFieldOutput    ClientShellField = "out"
	ClientShellFieldRemark    ClientShellField = "rem"
	ClientShellFieldTimeout   ClientShellField = "tmt"
	ClientShellFieldStatus    ClientShellField = "sta"
	ClientShellFieldCreatedAt ClientShellField = "c"
	ClientShellFieldUpdatedAt ClientShellField = "u"

	AllClientShellFields = []ClientShellField{
		ClientShellFieldInput,
		ClientShellFieldOutput,
		ClientShellFieldRemark,
		ClientShellFieldTimeout,
		ClientShellFieldStatus,
		ClientShellFieldCreatedAt,
		ClientShellFieldUpdatedAt,
	}
)

// client shell 的 dbkey
func clientShellFieldDBKey(clientId string, shellID string, shellField ClientShellField) []byte {
	return bytes.Join([][]byte{[]byte("client"), []byte(ClientFieldShell), []byte(clientId), []byte(shellField), []byte(shellID)}, []byte{'/'})
}
