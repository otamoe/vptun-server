package server

import (
	"bytes"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type (
	ClientShellSystem struct {
		clientSystem *ClientSystem
	}
)

func NewClientShellSystem(clientSystem *ClientSystem) (clientShellSystem *ClientShellSystem, err error) {
	clientShellSystem = &ClientShellSystem{
		clientSystem: clientSystem,
	}
	return
}
func (clientShellSystem *ClientShellSystem) List(clientId string, ltId string, limit int) (clientShells ClientShells, err error) {
	if !IsID(clientId) {
		err = ErrClientNotFound
		return
	}
	if ltId != "" && !IsID(ltId) {
		err = ErrClientShellNotFound
		return
	}
	err = clientShellSystem.clientSystem.db.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		seek := clientShellFieldDBKey(clientId, ltId, ClientShellFieldCreatedAt)
		prefix := clientShellFieldDBKey(clientId, "", ClientShellFieldCreatedAt)
		for it.Seek(seek); it.ValidForPrefix(prefix) && len(clientShells) < limit; it.Next() {
			item := it.Item()
			// key 相同
			if bytes.Equal(seek, item.Key()) {
				continue
			}

			keySplit := strings.Split(string(item.Key()), "/")
			id := string(keySplit[len(keySplit)-1])
			clientShell := newClientShell(clientId, id)
			if err = clientShell.load(txn); err != nil {
				return
			}
			clientShells = append(clientShells, clientShell)
		}
		return
	})
	if err != nil {
		clientShells = nil
	}
	return
}

func (clientShellSystem *ClientShellSystem) Get(clientId string, id string) (clientShell *ClientShell, err error) {
	if !IsID(clientId) {
		err = ErrClientNotFound
		return
	}
	if !IsID(id) {
		err = ErrClientShellNotFound
		return
	}
	err = clientShellSystem.clientSystem.db.View(func(txn *badger.Txn) (err error) {
		clientShell = newClientShell(clientId, id)
		if err = clientShell.load(txn); err != nil {
			return
		}
		return
	})
	if err != nil {
		clientShell = nil
	}
	return
}

func (clientShellSystem *ClientShellSystem) Delete(clientId string, id string) (clientShell *ClientShell, err error) {
	if !IsID(clientId) {
		err = ErrClientNotFound
		return
	}
	if !IsID(id) {
		err = ErrClientShellNotFound
		return
	}
	for i := 0; i < 512; i++ {
		err = clientShellSystem.clientSystem.db.Update(func(txn *badger.Txn) (err error) {
			clientShell = newClientShell(clientId, id)
			if err = clientShell.load(txn); err != nil {
				return
			}
			clientShell.delete(txn)
			return
		})
		// 无错误 响应
		if err == nil {
			break
		}
		// badger.ErrConflict 重试错误
		if err != badger.ErrConflict {
			break
		}
	}

	// 有错误
	if err != nil {
		clientShell = nil
	}
	return
}

func (clientShellSystem *ClientShellSystem) Save(clientId string, id string, cb func(clientShell *ClientShell) (rClientShell *ClientShell, err error)) (clientShell *ClientShell, err error) {
	if !IsID(clientId) {
		err = ErrClientNotFound
		return
	}
	if id != "" && !IsID(id) {
		err = ErrClientShellNotFound
		return
	}
	for i := 0; i < 512; i++ {
		if clientShellSystem.clientSystem.Get(clientId) == nil {
			err = ErrClientNotFound
			return
		}
		err = clientShellSystem.clientSystem.db.Update(func(txn *badger.Txn) (err error) {
			if id == "" {
				now := time.Now().UTC()
				clientShell = newClientShell(clientId, NewID(now)).
					WithStatus(-1).
					WithTimeout(600).
					WithCreatedAt(now).
					WithUpdatedAt(now)
			} else {
				clientShell = newClientShell(clientId, id)
				if err = clientShell.load(txn); err != nil {
					return
				}
			}

			// 回调修改
			if clientShell, err = cb(clientShell); err != nil {
				return
			}

			if id != "" {
				clientShell = clientShell.WithUpdatedAt(time.Now().UTC())
			}

			// 保存 client
			if err = clientShell.save(txn); err != nil {
				return
			}
			return
		})

		// 无错误 响应
		if err == nil {
			break
		}

		// badger.ErrConflict 重试错误
		if err != badger.ErrConflict {
			break
		}
	}

	// 有错误
	if err != nil {
		clientShell = nil
		return
	}

	return
}
