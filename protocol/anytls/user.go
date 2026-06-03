package anytls

import (
	"net"

	anytls "github.com/anytls/sing-anytls"
	"github.com/sagernet/sing-box/option"
)

func (h *Inbound) AddUsers(users []option.AnyTLSUser) error {
	for _, user := range users {
		h.uuidList = append(h.uuidList, user.Name)
	}
	userList := make([]anytls.User, len(h.uuidList))
	for i, uuid := range h.uuidList {
		userList[i] = anytls.User{Name: uuid, Password: uuid}
	}
	h.service.UpdateUsers(userList)
	return nil
}

func (h *Inbound) DelUsers(names []string) error {
	if len(names) == 0 {
		return nil
	}

	toDelete := make(map[string]struct{})
	for _, name := range names {
		toDelete[name] = struct{}{}
		h.userCons.Range(func(key, value interface{}) bool {
			if value.(string) == name {
				key.(net.Conn).Close()
				h.userCons.Delete(key)
			}
			return true
		})
	}

	remaining := make([]string, 0, len(h.uuidList))
	for _, uuid := range h.uuidList {
		if _, found := toDelete[uuid]; !found {
			remaining = append(remaining, uuid)
		}
	}

	h.uuidList = remaining
	userList := make([]anytls.User, len(h.uuidList))
	for i, uuid := range h.uuidList {
		userList[i] = anytls.User{Name: uuid, Password: uuid}
	}
	h.service.UpdateUsers(userList)
	return nil
}
