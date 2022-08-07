package fork

type MasterTool struct {
	f *Forker
}

func (m *MasterTool) ServerPipe() *ServerPipe {
	return m.f.serverPipe
}

func (m *MasterTool) RangeChildren(f func(cmd *Cmd) bool) {
	m.f.RangeChild(f)
}

func (m *MasterTool) SetForkPreRun(f ProcFunc) {
	m.f.SetPreForkChild(f)
}
