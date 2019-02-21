package utils

import (
	"log"
	"strings"
	"sync"

	xmpp "github.com/sshikaree/go-xmpp2"
)

// XMPPHub used to rule multiple XMPP connections
type XMPPHub struct {
	mu          sync.Mutex
	connections map[string]*xmpp.Client
}

// Register adds XMPP connection to hub
func (hub *XMPPHub) Register(c *xmpp.Client) {
	hub.mu.Lock()
	bare_jid := strings.Split(c.JID(), "/")[0]
	hub.connections[bare_jid] = c
	hub.mu.Unlock()
}

// Unregister removes XMPP connection from hub
func (hub *XMPPHub) Unregister(jid string) {
	hub.mu.Lock()
	hub.connections[jid].Close()
	delete(hub.connections, jid)
	hub.mu.Unlock()
}

// Len returns number of registered connections
func (hub *XMPPHub) Len() int {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	return len(hub.connections)
}

// SendBroadcast sends message to all registered connections
func (hub *XMPPHub) SendBroadcast(msg string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for _, c := range hub.connections {
		_, err := c.SendString(msg)
		if err != nil {
			log.Println(err)
			// c.Close()
			// delete(hub.connections, c)
		}
	}
}

// Range iterates over each connection and call f. If f returns false, it stops iteration.
func (hub *XMPPHub) Range(f func(jid string, client *xmpp.Client) bool) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for jid, client := range hub.connections {
		if !f(jid, client) {
			return
		}
	}
}

// Get connection by JID from hub
func (hub *XMPPHub) Get(jid string) (*xmpp.Client, bool) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	c, ok := hub.connections[jid]
	return c, ok
}

// NewXMPPHub creates new hub
func NewXMPPHub() *XMPPHub {
	hub := new(XMPPHub)
	hub.connections = make(map[string]*xmpp.Client)

	return hub
}
