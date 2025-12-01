package memory

const (
	MEMORY_SIZE uint32 = 0x10000
)

type Memory struct {
	memory [MEMORY_SIZE]uint8
}

func (m *Memory) Init() {
	m.memory = [MEMORY_SIZE]uint8{}
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.memory[addr]
}

func (m *Memory) Write(addr uint16, value uint8) {
	m.memory[addr] = value
}
