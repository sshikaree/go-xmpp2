//
// XEP-0030: Service Discovery
//

package xmpp

import (
	"encoding/xml"
	"fmt"
)

const (
	// IQTypeGet    = "get"
	// IQTypeSet    = "set"
	// IQTypeResult = "result"
	NSDiscoInfo  = "http://jabber.org/protocol/disco#info"
	NSDiscoItems = "http://jabber.org/protocol/disco#items"
)

type DiscoveryQuery struct {
	Items      []DiscoItem `xml:",omitempty"`
	Identities []Identity  `xml:",omitempty"`
	Features   []Feature   `xml:",omitempty"`
	Node       string      `xml:"node,attr,omitempty"`
}

type DiscoItem struct {
	XMLName xml.Name `xml:"item"`
	JID     string   `xml:"jid,attr"`
	Node    string   `xml:"node,attr,omitempty"`
	Name    string   `xml:"name,attr,omitempty"`
}

type Identity struct {
	XMLName  xml.Name `xml:"identity"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr,omitempty"`
}

type Feature struct {
	XMLName xml.Name `xml:"feature"`
	Var     string   `xml:"var,attr"`
}

// target is server name or room's JID or full JID with resource
func (c *Client) GetDiscoInfo(target, node string) error {
	var s string
	if node == "" {
		s = fmt.Sprintf(
			`<iq type='get'	from='%s' to='%s' id='info1'>
				<query xmlns='%s'/>
			</iq>`,
			c.JID(), target, NSDiscoInfo,
		)
	} else {
		s = fmt.Sprintf(
			`<iq type='get'	from='%s' to='%s' id='info1'>
				<query xmlns='%s' node='%s' />
			</iq>`,
			c.JID(), target, NSDiscoInfo, node,
		)
	}
	_, err := c.SendString(s)
	return err
}

// target is server name or room's JID or full JID with resource
func (c *Client) GetDiscoItems(target, node string) error {
	var s string
	if node == "" {
		s = fmt.Sprintf(
			`<iq type='get'	from='%s' to='%s' id='items1'>
				<query xmlns='%s'/>
			</iq>`,
			c.JID(), target, NSDiscoItems,
		)
	} else {
		s = fmt.Sprintf(
			`<iq type='get'	from='%s' to='%s' id='items1'>
				<query xmlns='%s' node='%s' />
			</iq>`,
			c.JID(), target, NSDiscoItems, node,
		)
	}
	_, err := c.SendString(s)
	return err
}

// // Returns request id
// func (c *Client) Discovery() (string, error) {
// 	const namespace = "http://jabber.org/protocol/disco#items"
// 	// use getCookie for a pseudo random id.
// 	reqID := strconv.FormatUint(uint64(getCookie()), 10)
// 	return reqID, c.RawInformationQuery(c.jid, c.domain, reqID, IQTypeGet, namespace, "")
// }

// // RawInformationQuery sends an information query request to the server.
// func (c *Client) RawInformationQuery(from, to, id, iqType, requestNamespace, body string) error {
// 	_, err := c.SendStrin(fmt.Sprintf(
// 		"<iq from='%s' to='%s' id='%s' type='%s'><query xmlns='%s'>%s</query></iq>",
// 		xmlEscape(from), xmlEscape(to), id, iqType, requestNamespace, body,
// 	))
// 	return err
// }
