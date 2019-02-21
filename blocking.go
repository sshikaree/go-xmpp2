//
// XEP-0191: Blocking Command
//

package xmpp

import "fmt"

// type BlockList struct {
// 	XMLName xml.Name        `xml:"urn:xmpp:blocking blocklist"`
// 	Items   []BlockListItem `xml:"item,omitempty"`
// }

const (
	NSBlock = "urn:xmpp:blocking"
)

type BlockList []BlockListItem

type BlockListItem struct {
	XMLName string `xml:"item"`
	JID     string `xml:"jid,attr"`
}

func (c *Client) Block(jid string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq from='%s' type='set' id='block1'>
			<block xmlns='%s'>
				<item jid='%s'/>
			</block>
		</iq>`,
		c.JID(), NSBlock, jid,
	))
	return err
}

func (c *Client) Unblock(jid string) error {
	_, err := c.SendString(fmt.Sprintf(
		`<iq type='set' id='unblock1'>
			<unblock xmlns='%s'>
				<item jid='%s'/>
			</unblock>
		</iq>`,
		NSBlock, jid,
	))

	return err
}

func (c *Client) GetBlocked() error {
	// simple version (maybe faster???)
	_, err := c.SendString(fmt.Sprintf(
		`<iq type='get' id='%s'>
			<blocklist xmlns='%s'/>
		  </iq>`,
		"blocklist1", NSBlock,
	))

	return err

	// iq := new(IQ)
	// iq.Type = "get"
	// iq.InnerElement = Query{
	// 	XMLName: xml.Name{
	// 		Local: "blocklist",
	// 		Space: NSBlock,
	// 	},
	// }
	// _, err := c.SendIQ(iq)
	// return err
}
