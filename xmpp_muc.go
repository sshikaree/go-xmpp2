// Copyright 2013 Flo Lauber <dev@qatfy.at>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(flo):
//   - support password protected MUC rooms
//   - cleanup signatures of join/leave functions

//
// XEP-0045: Multi-User Chat
//

package xmpp

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	nsMUC      = "http://jabber.org/protocol/muc"
	nsMUCUser  = "http://jabber.org/protocol/muc#user"
	nsMUCAdmin = "http://jabber.org/protocol/muc#admin"
	nsMUCOwner = "http://jabber.org/protocol/muc#owner"

	NoHistory      = 0
	CharHistory    = 1
	StanzaHistory  = 2
	SecondsHistory = 3
	SinceHistory   = 4
)

// Send sends room topic wrapped inside an XMPP message stanza body.
func (c *Client) SendTopic(room, topic string) (n int, err error) {
	// return c.SendString(fmt.Sprintf(
	// 	`<message to='%s' type='%s' xml:lang='en'><subject>%s</subject></message>`,
	// 	xmlEscape(room), "groupchat", xmlEscape(topic),
	// ))
	msg := new(Message)
	msg.To = room
	msg.Type = "groupchat"
	msg.Subject = topic

	return c.SendMessage(msg)
}

func (c *Client) JoinMUCNoHistory(room, nick string) (n int, err error) {
	if nick == "" {
		nick = strings.Split(c.jid, "@")[0]
	}
	return c.SendString(fmt.Sprintf(
		`<presence to='%s/%s'>
			<x xmlns='%s'>"
				<history maxchars='0'/>
			</x>
		</presence>`,
		xmlEscape(room), xmlEscape(nick), nsMUC,
	))
}

// xep-0045 7.2
func (c *Client) JoinMUC(room, nick string, history_type, history int, history_date *time.Time) (n int, err error) {
	if nick == "" {
		nick = strings.Split(c.jid, "@")[0]
	}
	switch history_type {
	case NoHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s' />
			</presence>`,
			xmlEscape(room), xmlEscape(nick), nsMUC,
		))
	case CharHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><history maxchars='%d'/></x>
			</presence>`,
			xmlEscape(room), xmlEscape(nick), nsMUC, history,
		))
	case StanzaHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><history maxstanzas='%d'/></x>
			</presence>`,
			xmlEscape(room), xmlEscape(nick), nsMUC, history,
		))
	case SecondsHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><history seconds='%d'/></x>
			</presence>`,
			xmlEscape(room), xmlEscape(nick), nsMUC, history,
		))
	case SinceHistory:
		if history_date != nil {
			return c.SendString(fmt.Sprintf(
				`<presence to='%s/%s'>
					<x xmlns='%s'><history since='%s'/></x>
				</presence>`,
				xmlEscape(room), xmlEscape(nick), nsMUC, history_date.Format(time.RFC3339),
			))
		}
	}
	return 0, errors.New("Unknown history option")
}

// xep-0045 7.2.6
func (c *Client) JoinProtectedMUC(jid, nick string, password string, history_type, history int, history_date *time.Time) (n int, err error) {
	if nick == "" {
		nick = strings.Split(c.jid, "@")[0]
	}
	switch history_type {
	case NoHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><password>%s</password></x>
			</presence>`,
			xmlEscape(jid), xmlEscape(nick), nsMUC, xmlEscape(password),
		))
	case CharHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><password>%s</password><history maxchars='%d'/></x>
			</presence>`,
			xmlEscape(jid), xmlEscape(nick), nsMUC, xmlEscape(password), history,
		))
	case StanzaHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><password>%s</password><history maxstanzas='%d'/></x>
			</presence>`,
			xmlEscape(jid), xmlEscape(nick), nsMUC, xmlEscape(password), history,
		))
	case SecondsHistory:
		return c.SendString(fmt.Sprintf(
			`<presence to='%s/%s'>
				<x xmlns='%s'><password>%s</password><history seconds='%d'/></x>
			</presence>`,
			xmlEscape(jid), xmlEscape(nick), nsMUC, xmlEscape(password), history,
		))
	case SinceHistory:
		if history_date != nil {
			return c.SendString(fmt.Sprintf(
				`<presence to='%s/%s'>
				<x xmlns='%s'><password>%s</password><history since='%s'/></x>
				</presence>`,
				xmlEscape(jid), xmlEscape(nick), nsMUC, xmlEscape(password), history_date.Format(time.RFC3339),
			))
		}
	}
	return 0, fmt.Errorf("Unknown history option: %d", history_type)
}

// xep-0045 7.14
func (c *Client) LeaveMUC(room string) (n int, err error) {
	return c.SendString(fmt.Sprintf(
		`<presence from='%s' to='%s' type='unavailable' />`,
		c.jid, xmlEscape(room),
	))
}

func (c *Client) SetRoomAffiliation(room, jid, affiliation, reason string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='ban1' to='%s' type='set'>
			<query xmlns='%s'>
				<item affiliation='%s' jid='%s'>
		  			<reason>%s</reason>
				</item>
			</query>		
		</iq>`,
		c.jid, xmlEscape(room), nsMUCAdmin, xmlEscape(affiliation), xmlEscape(jid), xmlEscape(reason),
	))
	return err
}

func (c *Client) Ban(room, jid, reason string) error {
	return c.SetRoomAffiliation(room, jid, "outcast", reason)
}

func (c *Client) ChangeNick(room, nick string) error {
	p := new(Presence)
	p.From = c.jid
	p.To = room + "/" + nick
	_, err := c.SendPresence(p)
	return err
}

type MUCForm struct {
}

func (c *Client) ConfigureRoom(room string, form MUCForm) error {
	return nil
}

// XEP-0249: Direct MUC Invitations
func (c *Client) DirectInvite(room, jid, password, reason string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<message from='%s' to='%s'>
			<x xmlns='jabber:x:conference'
				continue='true'
				jid='%s'
				password='%s'
				reason='%s'
				thread='e0ffe42b28561960c6b12b944a092794b9683a38'/>	
		</message>`,
		c.jid, xmlEscape(jid), xmlEscape(room), xmlEscape(password), xmlEscape(reason),
	))
	return err
}

// A user MAY have a reserved room nickname, for example through explicit room registration,
// database integration, or nickname "lockdown".
// A user SHOULD discover his or her reserved nickname before attempting to enter the room.
// This is done by sending a Service Discovery information request to the room JID
// while specifying a well-known Service Discovery node of "x-roomuser-item".
func (c *Client) DiscoverReservedNick(room string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='getnick1' to='%s' type='get'>
			<query xmlns='http://jabber.org/protocol/disco#info' node='x-roomuser-item'/>
		</iq>`,
		c.jid, xmlEscape(room),
	))
	return err
}

// XER-0045 10.2 Subsequent Room Configuration
func (c *Client) GetRoomConfig(room string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='ik3vs715' to='%s' type='get'>
			<query xmlns='http://jabber.org/protocol/muc#owner'/>
		</iq>`,
		c.jid, xmlEscape(room),
	))
	return err
}

func (c *Client) GetRoomMembers(room string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='member3'	to='%s' type='get'>
			<query xmlns='http://jabber.org/protocol/muc#admin'>
				<item affiliation='member'/>
			</query>
		</iq>`,
		c.jid, xmlEscape(room),
	))
	return err
}

func (c *Client) Invite(room string) error {
	return nil
}

func (c *Client) Kick(room, nick, reason string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='kick1' to='harfleur@chat.shakespeare.lit' type='set'>
			<query xmlns='http://jabber.org/protocol/muc#admin'>
				<item nick='%s' role='none'>
					<reason>%s</reason>
				</item>
			</query>
		</iq>`,
		c.jid, xmlEscape(nick), xmlEscape(reason),
	))
	return err
}

func (c *Client) RequestRoomVoice(room string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<message from='%s'	id='yd53c486' to='%s'>
			<x xmlns='jabber:x:data' type='submit'>
				<field var='FORM_TYPE'>
					<value>http://jabber.org/protocol/muc#request</value>
				</field>
				<field var='muc#role' type='list-single' label='Requested role'>
					<value>participant</value>
				</field>
			</x>
		</message>`,
		c.jid, xmlEscape(room),
	))
	return err
}

func (c *Client) SetRoomRole(room, nick, role, reason string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' id='voice1' to='%s' type='set'>
			<query xmlns='http://jabber.org/protocol/muc#admin'>
				<item nick='%s' role='%s'>
					<reason>%s</reason>
				</item>
			</query>
		</iq>`,
		c.jid, xmlEscape(room), xmlEscape(nick), xmlEscape(role), xmlEscape(reason),
	))
	return err
}
