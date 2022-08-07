package fork

type ChildrenTool struct {
	pipe    *ServerPipe
	cliPipe *ClientPipe
	childId int64
}

func (c *ChildrenTool) ServePipe() *ServerPipe {
	return c.pipe
}

func (c *ChildrenTool) ClientPipe() *ClientPipe {
	return c.cliPipe
}

func (c *ChildrenTool) ChildId() int64 {
	return c.childId
}
