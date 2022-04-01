package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"time"

	"github.com/dgraph-io/badger/v3"
	pb "github.com/otamoe/vptun-pb"
)

type (
	Client struct {
		*pb.Client

		iRouteAddress net.IP `json:"-"`
	}
	Clients []*Client
)

var ErrClientNotFound = &NotFoundError{
	Name: "client",
}

func (client *Client) Clone() (rClient *Client) {
	if client == nil {
		return nil
	}
	rClient = &Client{}
	*rClient = *client
	rClient.Client = &pb.Client{}

	if client.Client != nil {
		*rClient.Client = *client.Client
	}

	if client.iRouteAddress != nil {
		rClient.iRouteAddress = make(net.IP, len(client.iRouteAddress))
		copy(rClient.iRouteAddress, client.iRouteAddress)
	}

	if client.Client != nil && client.Client.Status != nil {
		rClient.Client.Status = &pb.Status{}
		*rClient.Client.Status = *client.Client.Status
	}
	return
}

func newClient(id string) *Client {
	return &Client{
		Client: &pb.Client{
			Id: id,
		},
	}
}

func (client *Client) WithKey(key string) (rClient *Client) {
	rClient = client.Clone()
	rClient.Key = key
	return rClient
}
func (client *Client) WithHostname(hostname string) (rClient *Client) {
	rClient = client.Clone()
	rClient.Hostname = hostname
	return rClient
}
func (client *Client) WithUserAgent(userAgent string) (rClient *Client) {
	rClient = client.Clone()
	rClient.UserAgent = userAgent
	return rClient
}
func (client *Client) WithConnectAddress(connectAddress string) (rClient *Client) {
	rClient = client.Clone()
	rClient.ConnectAddress = connectAddress
	return rClient
}

func (client *Client) WithRouteAddress(routeAddress net.IP) (rClient *Client) {
	rClient = client.Clone()
	rClient.iRouteAddress = routeAddress
	if len(routeAddress) == 0 {
		rClient.Client.RouteAddress = ""
	} else {
		rClient.Client.RouteAddress = routeAddress.String()
	}
	return rClient
}

func (client *Client) WithRemark(remark string) (rClient *Client) {
	rClient = client.Clone()
	rClient.Remark = remark
	return rClient
}
func (client *Client) WithShell(shell bool) (rClient *Client) {
	rClient = client.Clone()
	rClient.Shell = shell
	return rClient
}
func (client *Client) WithStatus(status *pb.Status) (rClient *Client) {
	rClient = client.Clone()
	rClient.Status = status
	return rClient
}
func (client *Client) WithState(state pb.State) (rClient *Client) {
	rClient = client.Clone()
	rClient.State = state
	return rClient
}
func (client *Client) WithOnline(online bool) (rClient *Client) {
	rClient = client.Clone()
	rClient.Online = online
	return rClient
}
func (client *Client) WithCreatedAt(t time.Time) (rClient *Client) {
	rClient = client.Clone()
	rClient.CreatedAt = t.Unix()
	return rClient
}
func (client *Client) WithUpdatedAt(t time.Time) (rClient *Client) {
	rClient = client.Clone()
	rClient.UpdatedAt = t.Unix()
	return rClient
}
func (client *Client) WithConnectAt(t time.Time) (rClient *Client) {
	rClient = client.Clone()
	rClient.ConnectAt = t.Unix()
	return rClient
}
func (client *Client) WithExpiredAt(t time.Time) (rClient *Client) {
	rClient = client.Clone()
	rClient.ExpiredAt = t.Unix()
	return rClient
}

func (client *Client) load(txn *badger.Txn) (err error) {
	var item *badger.Item
	var value []byte

	// ClientFieldKey
	{
		client.Client.Key = ""
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldKey)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Client.Key = string(value)
			}
		}
	}
	// ClientFieldHostname
	{
		client.Client.Hostname = ""
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldHostname)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Client.Hostname = string(value)
			}
		}
	}
	// ClientFieldUserAgent
	{
		client.Client.UserAgent = ""
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldUserAgent)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Client.UserAgent = string(value)
			}
		}
	}
	// ClientFieldConnectAddress
	{
		client.Client.ConnectAddress = ""
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldConnectAddress)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Client.ConnectAddress = string(value)
			}
		}
	}

	// ClientFieldRouteAddress
	{
		client.Client.RouteAddress = ""
		client.iRouteAddress = nil
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldRouteAddress)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Client.RouteAddress = string(value)
				client.iRouteAddress = net.ParseIP(client.Client.RouteAddress)
				if len(client.RouteAddress) == 0 {
					client.Client.RouteAddress = ""
				}
			}
		}
	}
	// ClientFieldRemark
	{
		client.Client.Remark = ""
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldRemark)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Remark = string(value)
			}
		}
	}

	// ClientFieldOnline
	{
		client.Online = false
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldOnline)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 && value[0] == 1 {
				client.Online = true
			}
		}
	}
	// ClientFieldShell
	{
		client.Shell = false
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldShell)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 && value[0] == 1 {
				client.Shell = true
			}
		}
	}
	// ClientFieldStatus
	{
		client.Status = nil
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldStatus)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				client.Status = &pb.Status{}
				if err = client.Status.Unmarshal(value); err != nil {
					return
				}
			}
		}
	}

	// ClientFieldState
	{
		client.State = 0
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldState)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &client.State)
			}
		}
	}

	// ClientFieldCreatedAt
	{
		client.CreatedAt = 0
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldCreatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}

		// 创建时间是必须的否则未找到
		if err == badger.ErrKeyNotFound {
			err = ErrClientNotFound
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &client.CreatedAt)
			}
		}
	}
	// ClientFieldUpdatedAt
	{
		client.UpdatedAt = 0
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldUpdatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &client.UpdatedAt)
			}
		}
	}
	// ClientFieldConnectAt
	{
		client.ConnectAt = 0
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldConnectAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &client.ConnectAt)
			}
		}
	}

	// ClientFieldExpiredAt
	{
		client.ExpiredAt = 0
		if item, err = txn.Get(clientFieldDBKey(client.Id, ClientFieldExpiredAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &client.ExpiredAt)
			}
		}
	}

	err = nil
	return
}

func (client *Client) UnmarshalJSON(data []byte) (err error) {
	client.Client = &pb.Client{}
	if err = json.Unmarshal(data, client.Client); err != nil {
		return
	}

	// RouteAddress 必须 是 ip
	if client.Client.RouteAddress != "" {
		if client.iRouteAddress = net.ParseIP(client.Client.RouteAddress); client.iRouteAddress == nil {
			err = &ValidateError{
				Name: "routeAddress",
			}
			return
		}
	}
	return
}

func (client *Client) save(txn *badger.Txn) (err error) {
	if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldKey), []byte(client.Client.Key)); err != nil {
		return
	}
	if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldHostname), []byte(client.Client.Hostname)); err != nil {
		return
	}
	if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldConnectAddress), []byte(client.Client.ConnectAddress)); err != nil {
		return
	}
	if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldRouteAddress), []byte(client.Client.RouteAddress)); err != nil {
		return
	}
	if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldRemark), []byte(client.Client.Remark)); err != nil {
		return
	}

	{
		v := byte(0)
		if client.Client.Online {
			v = 1
		}
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldOnline), []byte{v}); err != nil {
			return
		}
	}
	{
		v := byte(0)
		if client.Client.Shell {
			v = 1
		}
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldShell), []byte{v}); err != nil {
			return
		}
	}
	{
		if client.Client.Status == nil {
			if err = txn.Delete(clientFieldDBKey(client.Id, ClientFieldStatus)); err != nil {
				return
			}
		} else {
			var v []byte
			if v, err = client.Client.Status.Marshal(); err != nil {
				return
			}
			if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldStatus), v); err != nil {
				return
			}
		}
	}

	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, client.State)
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldState), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, client.CreatedAt)
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldCreatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, client.UpdatedAt)
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldUpdatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, client.ConnectAt)
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldConnectAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, client.ExpiredAt)
		if err = txn.Set(clientFieldDBKey(client.Id, ClientFieldExpiredAt), buf.Bytes()); err != nil {
			return
		}
	}
	return
}

func (client *Client) delete(txn *badger.Txn) (err error) {
	for _, field := range AllClientFields {
		if err = txn.Delete(clientFieldDBKey(client.Id, field)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
	}

	// 删除 shell 所有数据
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()
	prefix := clientFieldDBKey(client.Id, ClientFieldShell)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		if err = txn.Delete(item.KeyCopy(nil)); err != nil {
			return
		}
	}
	return
}

func (clients Clients) Len() int {
	return len(clients)
}
func (clients Clients) Less(i, j int) bool {
	return bytes.Compare(clients[i].iRouteAddress, clients[j].iRouteAddress) != 1
}

func (clients Clients) Swap(i, j int) {
	clients[i], clients[j] = clients[j], clients[i]
}
