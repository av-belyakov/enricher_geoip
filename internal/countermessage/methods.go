package countermessage

// NewSomeMessage некоторое новое сообщение счетчика
func NewSomeMessage() *SomeMessage {
	return &SomeMessage{}
}

func (c *SomeMessage) GetType() string {
	return c.Type
}

func (c *SomeMessage) SetType(v string) {
	c.Type = v
}

func (c *SomeMessage) GetMessage() string {
	return c.Message
}

func (c *SomeMessage) SetMessage(v string) {
	c.Message = v
}
