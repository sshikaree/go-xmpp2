package xmpp

import (
	"encoding/xml"
	"fmt"
)

const (
	NSRoster = "jabber:iq:roster"
)

// Roster asks for the chat roster.
func (c *Client) Roster() error {
	_, err := c.SendString(fmt.Sprintf(
		"<iq from='%s' type='get' id='roster1'><query xmlns='jabber:iq:roster'/></iq>\n",
		xmlEscape(c.jid),
	))
	return err
}

func (c *Client) UpdateRosterItem(item *RosterItem) error {
	roster := Roster{}
	roster = append(roster, *item)

	rosterBytes, err := xml.Marshal(&roster)
	if err != nil {
		return err
	}
	iq := new(IQ)
	iq.Type = "set"
	iq.InnerElement = Query{
		XMLName: xml.Name{
			Space: NSRoster,
			Local: "query",
		},
		InnerXML: rosterBytes,
	}
	_, err = c.SendIQ(iq)

	return err
}

func (c *Client) RemoveRosterItem(jid string) error {
	item := new(RosterItem)
	item.JID = jid
	item.Subscription = "remove"
	return c.UpdateRosterItem(item)
}

func (c *Client) ApproveSubscription(jid string) error {
	_, err := c.SendString(fmt.Sprintf(
		"<presence to='%s' type='subscribed'/>",
		xmlEscape(jid),
	))
	return err
}

func (c *Client) RevokeSubscription(jid string) error {
	_, err := c.SendString(fmt.Sprintf(
		"<presence to='%s' type='unsubscribed'/>",
		xmlEscape(jid),
	))
	return err
}

func (c *Client) RequestSubscription(jid string) error {
	_, err := c.SendString(fmt.Sprintf(
		"<presence to='%s' type='subscribe'/>",
		xmlEscape(jid),
	))
	return err
}
