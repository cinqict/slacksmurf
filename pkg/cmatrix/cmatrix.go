package cmatrix

import (
	"fmt"

	"github.com/nlopes/slack"
)

//Handler stores a Prefix and a channel
type Handler struct {
	Prefix         string
	HandlerChannel chan HandlerChannel
}

//Handlers is a collection of multiple handler objects
type Handlers []Handler

//HandlerChannel is a communication channel between our brain and the plugin
type HandlerChannel struct {
	Channel *slack.Channel
	Event   *slack.MessageEvent
	UserID  string
}

//AttachmentChannel represents a channel that handles returns message to slack
type AttachmentChannel struct {
	Channel      *slack.Channel
	Attachment   *slack.Attachment
	DisplayTitle string
}

//Cmatrix is our main object
type Cmatrix struct {
	Handlers
	returnChannel chan AttachmentChannel
}

var c *Cmatrix

func init() {
	c = New()
	rc := make(chan AttachmentChannel)
	c.returnChannel = rc
}

//Errors

//PrefixAlreadyExists denotes encountering an already existing prefix
type PrefixAlreadyExists string

// Error returns the formatted prefix error.
func (str PrefixAlreadyExists) Error() string {
	return fmt.Sprintf("Prefix Already exists in CMatrix %q", string(str))
}

//ChannelNotFound denotes encountering a request for a unknown channel based on
// prefix
type ChannelNotFound string

//Error returns the formatted ChannelNotFound error.
func (str ChannelNotFound) Error() string {
	return fmt.Sprintf("Channel beloning to prefix not found in CMatrix %q", string(str))
}

//New instantiate a new Cmatrix instance
func New() *Cmatrix {
	c = new(Cmatrix)
	return c
}

//Add add a handler to the singleton cmatrix instance
func Add(prefix string, ch chan HandlerChannel) error {
	return c.add(prefix, ch)
}

//add private method to handle the add
func (c *Cmatrix) add(prefix string, ch chan HandlerChannel) error {
	//check if the prefix already exists .. if it does return error
	if c.prefixExists(prefix) {
		return PrefixAlreadyExists(prefix)
	}

	//compose the handler
	h := Handler{
		Prefix:         prefix,
		HandlerChannel: ch,
	}

	//append the handler to the slice of handlers
	c.Handlers = append(c.Handlers, h)

	return nil
}

//CGetByP a channel by prefix
func CGetByP(prefix string) (chan HandlerChannel, error) {
	return c.cGetByP(prefix)
}

//cGetByP searches Cmatrix for channel
func (c *Cmatrix) cGetByP(prefix string) (chan HandlerChannel, error) {
	for _, h := range c.Handlers {
		if h.Prefix == prefix {
			return h.HandlerChannel, nil
		}
	}
	return nil, ChannelNotFound(prefix)
}

//GetReturnChannel returns the slack return channel
func GetReturnChannel() chan AttachmentChannel {
	return c.getReturnChannel()
}

func (c *Cmatrix) getReturnChannel() chan AttachmentChannel {
	return c.returnChannel
}

//return true if the prefix already exists in the Cmatrix
func (c *Cmatrix) prefixExists(prefix string) bool {
	for _, h := range c.Handlers {
		if h.Prefix == prefix {
			return true
		}
	}
	return false
}
