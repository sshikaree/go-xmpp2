//XEP-0047: In-Band Bytestreams

package xmpp

import (
	"encoding/xml"
)

const (
	NSIBB = "http://jabber.org/protocol/ibb"
)

type IBBOpen struct {
	XMLName   xml.Name `xml:"http://jabber.org/protocol/ibb open"`
	BlockSize int      `xml:"block-size,attr"`
	SID       string   `xml:"sid,attr"`
	Stanza    string   `xml:"stanza,attr,omitempty"`
}

type IBBData struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/ibb data"`
	Seq     int      `xml:"seq,attr"`
	SID     string   `xml:"sid,attr"`
	Payload []byte   `xml:",chardata"`
}

type IBBCLose struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/ibb close"`
	SID     string   `xml:"sid,attr"`
}
