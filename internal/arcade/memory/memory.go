package memory

type Memory struct {
	memory [10000]uint8
}

func (m *Memory) Init() {
	m.memory = [10000]uint8{}
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.memory[addr]
}

func (m *Memory) Write(addr uint16, value uint8) {
	m.memory[addr] = value
}
