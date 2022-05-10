package server

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type (
	ClientShell struct {
		ClientId  string `json:"clientId,omitempty"`
		Id        string `json:"id,omitempty"`
		Input     string `json:"input,omitempty"`
		Output    string `json:"output,omitempty"`
		Remark    string `json:"remark,omitempty"`
		Timeout   uint32 `json:"timeout,omitempty"`
		Status    int32  `json:"status,omitempty"`
		CreatedAt int64  `json:"createdAt,omitempty"`
		UpdatedAt int64  `json:"updatedAt,omitempty"`
	}
	ClientShells []*ClientShell
)

var (
	ErrClientShellNotFound = &NotFoundError{
		Name: "client shell",
	}
)

func (clientShell ClientShell) Clone() (rClientShell *ClientShell) {
	rClientShell = &ClientShell{}
	*rClientShell = clientShell
	return
}

func (clientShell *ClientShell) WithInput(input string) (rClientShell *ClientShell) {
	rClientShell = clientShell.Clone()
	rClientShell.Input = input
	return rClientShell
}
func (clientShell *ClientShell) WithOutput(input string) (rClientShell *ClientShell) {
	rClientShell = clientShell.Clone()
	rClientShell.Output = input
	return rClientShell
}

func (clientShell *ClientShell) WithRemark(remark string) (rClientShell *ClientShell) {
	rClientShell = clientShell.Clone()
	rClientShell.Remark = remark
	return rClientShell
}
func (clientShell *ClientShell) WithTimeout(timeout uint32) (rClientShell *ClientShell) {
	rClientShell = clientShell.Clone()
	rClientShell.Timeout = timeout
	return rClientShell
}
func (clientShell *ClientShell) WithStatus(status int32) (rClientShell *ClientShell) {
	rClientShell = clientShell.Clone()
	rClientShell.Status = status
	return rClientShell
}
func (client *ClientShell) WithCreatedAt(t time.Time) (rClientShell *ClientShell) {
	rClientShell = client.Clone()
	rClientShell.CreatedAt = t.Unix()
	return rClientShell
}
func (client *ClientShell) WithUpdatedAt(t time.Time) (rClientShell *ClientShell) {
	rClientShell = client.Clone()
	rClientShell.UpdatedAt = t.Unix()
	return rClientShell
}

func (clientShell *ClientShell) load(txn *badger.Txn) (err error) {
	var item *badger.Item
	var value []byte

	// ClientShellFieldInput
	{
		clientShell.Input = ""
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldInput)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				clientShell.Input = string(value)
			}
		}
	}
	// ClientShellFieldOutput
	{
		clientShell.Output = ""
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldOutput)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				clientShell.Output = string(value)
			}
		}
	}
	// ClientShellFieldRemark
	{
		clientShell.Remark = ""
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldRemark)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				clientShell.Remark = string(value)
			}
		}
	}

	// ClientShellFieldTimeout
	{
		clientShell.Timeout = 0
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldTimeout)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &clientShell.Timeout)
			}
		}
	}

	// ClientShellFieldStatus
	{
		clientShell.Status = 0
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldStatus)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &clientShell.Status)
			}
		}
	}

	// ClientShellFieldCreatedAt
	{
		clientShell.CreatedAt = 0
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldCreatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}

		// 创建时间是必须的否则未找到
		if err == badger.ErrKeyNotFound {
			err = ErrClientShellNotFound
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &clientShell.CreatedAt)
			}
		}
	}
	// ClientShellFieldUpdatedAt
	{
		clientShell.UpdatedAt = 0
		if item, err = txn.Get(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldUpdatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &clientShell.UpdatedAt)
			}
		}
	}

	return
}

func (clientShell *ClientShell) save(txn *badger.Txn) (err error) {
	if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldInput), []byte(clientShell.Input)); err != nil {
		return
	}
	if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldOutput), []byte(clientShell.Output)); err != nil {
		return
	}
	if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldRemark), []byte(clientShell.Remark)); err != nil {
		return
	}

	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, clientShell.Timeout)
		if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldTimeout), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, clientShell.Status)
		if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldStatus), buf.Bytes()); err != nil {
			return
		}
	}

	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, clientShell.CreatedAt)
		if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldCreatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, clientShell.UpdatedAt)
		if err = txn.Set(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, ClientShellFieldUpdatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	return
}

func (clientShell *ClientShell) delete(txn *badger.Txn) (err error) {
	for _, field := range AllClientShellFields {
		if err = txn.Delete(clientShellFieldDBKey(clientShell.ClientId, clientShell.Id, field)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
	}
	return
}

func (clientShells ClientShells) Len() int {
	return len(clientShells)
}
func (clientShells ClientShells) Less(i, j int) bool {
	return strings.Compare(clientShells[i].Id, clientShells[j].Id) != 1
}

func (clientShells ClientShells) Swap(i, j int) {
	clientShells[i], clientShells[j] = clientShells[j], clientShells[i]
}

func newClientShell(clientId string, id string) *ClientShell {
	return &ClientShell{
		Id:       id,
		ClientId: clientId,
	}
}
