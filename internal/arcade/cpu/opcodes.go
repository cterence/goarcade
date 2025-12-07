package cpu

type inst struct {
	Name   string
	Op1    string
	Op2    string
	Length uint8
	Cycles uint8
	exec   func(*CPU, string)
}

var InstByOpcode = [256]inst{
	// 0x00-0x0F
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},     // 0x00
	{Name: "LXI", Op1: "BC", Op2: "", Length: 3, Cycles: 10, exec: lxi},  // 0x01
	{Name: "STAX", Op1: "BC", Op2: "", Length: 1, Cycles: 7, exec: stax}, // 0x02
	{Name: "INX", Op1: "BC", Op2: "", Length: 1, Cycles: 5, exec: inx},   // 0x03
	{Name: "INR", Op1: "B", Op2: "", Length: 1, Cycles: 5, exec: inr},    // 0x04
	{Name: "DCR", Op1: "B", Op2: "", Length: 1, Cycles: 5, exec: dcr},    // 0x05
	{Name: "MVI", Op1: "B", Op2: "", Length: 2, Cycles: 7, exec: mvi},    // 0x06
	{Name: "RLC", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: rlc},     // 0x07
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},     // 0x08 *NOP
	{Name: "DAD", Op1: "BC", Op2: "", Length: 1, Cycles: 10, exec: dad},  // 0x09
	{Name: "LDAX", Op1: "BC", Op2: "", Length: 1, Cycles: 7, exec: ldax}, // 0x0A
	{Name: "DCX", Op1: "BC", Op2: "", Length: 1, Cycles: 5, exec: dcx},   // 0x0B
	{Name: "INR", Op1: "C", Op2: "", Length: 1, Cycles: 5, exec: inr},    // 0x0C
	{Name: "DCR", Op1: "C", Op2: "", Length: 1, Cycles: 5, exec: dcr},    // 0x0D
	{Name: "MVI", Op1: "C", Op2: "", Length: 2, Cycles: 7, exec: mvi},    // 0x0E
	{Name: "RRC", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: rrc},     // 0x0F

	// 0x10-0x1F
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},     // 0x10 *NOP
	{Name: "LXI", Op1: "DE", Op2: "", Length: 3, Cycles: 10, exec: lxi},  // 0x11
	{Name: "STAX", Op1: "DE", Op2: "", Length: 1, Cycles: 7, exec: stax}, // 0x12
	{Name: "INX", Op1: "DE", Op2: "", Length: 1, Cycles: 5, exec: inx},   // 0x13
	{Name: "INR", Op1: "D", Op2: "", Length: 1, Cycles: 5, exec: inr},    // 0x14
	{Name: "DCR", Op1: "D", Op2: "", Length: 1, Cycles: 5, exec: dcr},    // 0x15
	{Name: "MVI", Op1: "D", Op2: "", Length: 2, Cycles: 7, exec: mvi},    // 0x16
	{Name: "RAL", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: ral},     // 0x17
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},     // 0x18 *NOP
	{Name: "DAD", Op1: "DE", Op2: "", Length: 1, Cycles: 10, exec: dad},  // 0x19
	{Name: "LDAX", Op1: "DE", Op2: "", Length: 1, Cycles: 7, exec: ldax}, // 0x1A
	{Name: "DCX", Op1: "DE", Op2: "", Length: 1, Cycles: 5, exec: dcx},   // 0x1B
	{Name: "INR", Op1: "E", Op2: "", Length: 1, Cycles: 5, exec: inr},    // 0x1C
	{Name: "DCR", Op1: "E", Op2: "", Length: 1, Cycles: 5, exec: dcr},    // 0x1D
	{Name: "MVI", Op1: "E", Op2: "", Length: 2, Cycles: 7, exec: mvi},    // 0x1E
	{Name: "RAR", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: rar},     // 0x1F

	// 0x20-0x2F
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},                                         // 0x20 *NOP
	{Name: "LXI", Op1: "HL", Op2: "", Length: 3, Cycles: 10, exec: lxi},                                      // 0x21
	{Name: "SHLD", Op1: "", Op2: "", Length: 3, Cycles: 16, exec: shld},                                      // 0x22
	{Name: "INX", Op1: "HL", Op2: "", Length: 1, Cycles: 5, exec: inx},                                       // 0x23
	{Name: "INR", Op1: "H", Op2: "", Length: 1, Cycles: 5, exec: inr},                                        // 0x24
	{Name: "DCR", Op1: "H", Op2: "", Length: 1, Cycles: 5, exec: dcr},                                        // 0x25
	{Name: "MVI", Op1: "H", Op2: "", Length: 2, Cycles: 7, exec: mvi},                                        // 0x26
	{Name: "DAA", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: daa},                                         // 0x27
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},                                         // 0x28 *NOP
	{Name: "DAD", Op1: "HL", Op2: "", Length: 1, Cycles: 10, exec: dad},                                      // 0x29
	{Name: "LHLD", Op1: "", Op2: "", Length: 3, Cycles: 16, exec: lhld},                                      // 0x2A
	{Name: "DCX", Op1: "HL", Op2: "", Length: 1, Cycles: 5, exec: dcx},                                       // 0x2B
	{Name: "INR", Op1: "L", Op2: "", Length: 1, Cycles: 5, exec: inr},                                        // 0x2C
	{Name: "DCR", Op1: "L", Op2: "", Length: 1, Cycles: 5, exec: dcr},                                        // 0x2D
	{Name: "MVI", Op1: "L", Op2: "", Length: 2, Cycles: 7, exec: mvi},                                        // 0x2E
	{Name: "CMA", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.A = 0xFF - c.A }}, // 0x2F

	// 0x30-0x3F
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},                                                  // 0x30 *NOP
	{Name: "LXI", Op1: "SP", Op2: "", Length: 3, Cycles: 10, exec: lxi},                                               // 0x31
	{Name: "STA", Op1: "", Op2: "", Length: 3, Cycles: 13, exec: sta},                                                 // 0x32
	{Name: "INX", Op1: "SP", Op2: "", Length: 1, Cycles: 5, exec: inx},                                                // 0x33
	{Name: "INR", Op1: "M", Op2: "", Length: 1, Cycles: 10, exec: inr},                                                // 0x34
	{Name: "DCR", Op1: "M", Op2: "", Length: 1, Cycles: 10, exec: dcr},                                                // 0x35
	{Name: "MVI", Op1: "M", Op2: "", Length: 2, Cycles: 10, exec: mvi},                                                // 0x36
	{Name: "STC", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.setCYF(true) }},            // 0x37
	{Name: "NOP", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: nop},                                                  // 0x38 *NOP
	{Name: "DAD", Op1: "SP", Op2: "", Length: 1, Cycles: 10, exec: dad},                                               // 0x39
	{Name: "LDA", Op1: "", Op2: "", Length: 3, Cycles: 13, exec: lda},                                                 // 0x3A
	{Name: "DCX", Op1: "SP", Op2: "", Length: 1, Cycles: 5, exec: dcx},                                                // 0x3B
	{Name: "INR", Op1: "A", Op2: "", Length: 1, Cycles: 5, exec: inr},                                                 // 0x3C
	{Name: "DCR", Op1: "A", Op2: "", Length: 1, Cycles: 5, exec: dcr},                                                 // 0x3D
	{Name: "MVI", Op1: "A", Op2: "", Length: 2, Cycles: 7, exec: mvi},                                                 // 0x3E
	{Name: "CMC", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.setCYF(c.getCYF() == 0) }}, // 0x3F

	// 0x40-0x4F
	{Name: "MOV", Op1: "B", Op2: "B", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x40
	{Name: "MOV", Op1: "B", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.C }},                                      // 0x41
	{Name: "MOV", Op1: "B", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.D }},                                      // 0x42
	{Name: "MOV", Op1: "B", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.E }},                                      // 0x43
	{Name: "MOV", Op1: "B", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.H }},                                      // 0x44
	{Name: "MOV", Op1: "B", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.L }},                                      // 0x45
	{Name: "MOV", Op1: "B", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.B = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x46
	{Name: "MOV", Op1: "B", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.B = c.A }},                                      // 0x47
	{Name: "MOV", Op1: "C", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.B }},                                      // 0x48
	{Name: "MOV", Op1: "C", Op2: "C", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x49
	{Name: "MOV", Op1: "C", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.D }},                                      // 0x4A
	{Name: "MOV", Op1: "C", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.E }},                                      // 0x4B
	{Name: "MOV", Op1: "C", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.H }},                                      // 0x4C
	{Name: "MOV", Op1: "C", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.L }},                                      // 0x4D
	{Name: "MOV", Op1: "C", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.C = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x4E
	{Name: "MOV", Op1: "C", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.C = c.A }},                                      // 0x4F

	// 0x50-0x5F
	{Name: "MOV", Op1: "D", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.B }},                                      // 0x50
	{Name: "MOV", Op1: "D", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.C }},                                      // 0x51
	{Name: "MOV", Op1: "D", Op2: "D", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x52
	{Name: "MOV", Op1: "D", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.E }},                                      // 0x53
	{Name: "MOV", Op1: "D", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.H }},                                      // 0x54
	{Name: "MOV", Op1: "D", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.L }},                                      // 0x55
	{Name: "MOV", Op1: "D", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.D = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x56
	{Name: "MOV", Op1: "D", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.D = c.A }},                                      // 0x57
	{Name: "MOV", Op1: "E", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.B }},                                      // 0x58
	{Name: "MOV", Op1: "E", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.C }},                                      // 0x59
	{Name: "MOV", Op1: "E", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.D }},                                      // 0x5A
	{Name: "MOV", Op1: "E", Op2: "E", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x5B
	{Name: "MOV", Op1: "E", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.H }},                                      // 0x5C
	{Name: "MOV", Op1: "E", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.L }},                                      // 0x5D
	{Name: "MOV", Op1: "E", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.E = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x5E
	{Name: "MOV", Op1: "E", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.E = c.A }},                                      // 0x5F

	// 0x60-0x6F
	{Name: "MOV", Op1: "H", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.B }},                                      // 0x60
	{Name: "MOV", Op1: "H", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.C }},                                      // 0x61
	{Name: "MOV", Op1: "H", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.D }},                                      // 0x62
	{Name: "MOV", Op1: "H", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.E }},                                      // 0x63
	{Name: "MOV", Op1: "H", Op2: "H", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x64
	{Name: "MOV", Op1: "H", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.L }},                                      // 0x65
	{Name: "MOV", Op1: "H", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.H = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x66
	{Name: "MOV", Op1: "H", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.H = c.A }},                                      // 0x67
	{Name: "MOV", Op1: "L", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.B }},                                      // 0x68
	{Name: "MOV", Op1: "L", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.C }},                                      // 0x69
	{Name: "MOV", Op1: "L", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.D }},                                      // 0x6A
	{Name: "MOV", Op1: "L", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.E }},                                      // 0x6B
	{Name: "MOV", Op1: "L", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.H }},                                      // 0x6C
	{Name: "MOV", Op1: "L", Op2: "L", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x6D
	{Name: "MOV", Op1: "L", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.L = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x6E
	{Name: "MOV", Op1: "L", Op2: "A", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.L = c.A }},                                      // 0x6F

	// 0x70-0x7F
	{Name: "MOV", Op1: "M", Op2: "B", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.B) }},   // 0x70
	{Name: "MOV", Op1: "M", Op2: "C", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.C) }},   // 0x71
	{Name: "MOV", Op1: "M", Op2: "D", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.D) }},   // 0x72
	{Name: "MOV", Op1: "M", Op2: "E", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.E) }},   // 0x73
	{Name: "MOV", Op1: "M", Op2: "H", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.H) }},   // 0x74
	{Name: "MOV", Op1: "M", Op2: "L", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.L) }},   // 0x75
	{Name: "HLT", Op1: "", Op2: "", Length: 1, Cycles: 7, exec: hlt},                                                                         // 0x76
	{Name: "MOV", Op1: "M", Op2: "A", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.Bus.Write(uint16(c.H)<<8|uint16(c.L), c.A) }},   // 0x77
	{Name: "MOV", Op1: "A", Op2: "B", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.B }},                                      // 0x78
	{Name: "MOV", Op1: "A", Op2: "C", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.C }},                                      // 0x79
	{Name: "MOV", Op1: "A", Op2: "D", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.D }},                                      // 0x7A
	{Name: "MOV", Op1: "A", Op2: "E", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.E }},                                      // 0x7B
	{Name: "MOV", Op1: "A", Op2: "H", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.H }},                                      // 0x7C
	{Name: "MOV", Op1: "A", Op2: "L", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.A = c.L }},                                      // 0x7D
	{Name: "MOV", Op1: "A", Op2: "M", Length: 1, Cycles: 7, exec: func(c *CPU, _ string) { c.A = c.Bus.Read(uint16(c.H)<<8 | uint16(c.L)) }}, // 0x7E
	{Name: "MOV", Op1: "A", Op2: "A", Length: 1, Cycles: 5, exec: nop},                                                                       // 0x7F

	// 0x80-0x8F
	{Name: "ADD", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x80
	{Name: "ADD", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x81
	{Name: "ADD", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x82
	{Name: "ADD", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x83
	{Name: "ADD", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x84
	{Name: "ADD", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x85
	{Name: "ADD", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: add}, // 0x86
	{Name: "ADD", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: add}, // 0x87
	{Name: "ADC", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x88
	{Name: "ADC", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x89
	{Name: "ADC", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x8A
	{Name: "ADC", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x8B
	{Name: "ADC", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x8C
	{Name: "ADC", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x8D
	{Name: "ADC", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: adc}, // 0x8E
	{Name: "ADC", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: adc}, // 0x8F

	// 0x90-0x9F
	{Name: "SUB", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x90
	{Name: "SUB", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x91
	{Name: "SUB", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x92
	{Name: "SUB", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x93
	{Name: "SUB", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x94
	{Name: "SUB", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x95
	{Name: "SUB", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: sub}, // 0x96
	{Name: "SUB", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: sub}, // 0x97
	{Name: "SBB", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x98
	{Name: "SBB", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x99
	{Name: "SBB", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x9A
	{Name: "SBB", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x9B
	{Name: "SBB", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x9C
	{Name: "SBB", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x9D
	{Name: "SBB", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: sbb}, // 0x9E
	{Name: "SBB", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: sbb}, // 0x9F

	// 0xA0-0xAF
	{Name: "ANA", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA0
	{Name: "ANA", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA1
	{Name: "ANA", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA2
	{Name: "ANA", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA3
	{Name: "ANA", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA4
	{Name: "ANA", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA5
	{Name: "ANA", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: ana}, // 0xA6
	{Name: "ANA", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: ana}, // 0xA7
	{Name: "XRA", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xA8
	{Name: "XRA", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xA9
	{Name: "XRA", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xAA
	{Name: "XRA", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xAB
	{Name: "XRA", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xAC
	{Name: "XRA", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xAD
	{Name: "XRA", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: xra}, // 0xAE
	{Name: "XRA", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: xra}, // 0xAF

	// 0xB0-0xBF
	{Name: "ORA", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB0
	{Name: "ORA", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB1
	{Name: "ORA", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB2
	{Name: "ORA", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB3
	{Name: "ORA", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB4
	{Name: "ORA", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB5
	{Name: "ORA", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: ora}, // 0xB6
	{Name: "ORA", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: ora}, // 0xB7
	{Name: "CMP", Op1: "B", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xB8
	{Name: "CMP", Op1: "C", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xB9
	{Name: "CMP", Op1: "D", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xBA
	{Name: "CMP", Op1: "E", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xBB
	{Name: "CMP", Op1: "H", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xBC
	{Name: "CMP", Op1: "L", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xBD
	{Name: "CMP", Op1: "M", Op2: "", Length: 1, Cycles: 7, exec: cmp}, // 0xBE
	{Name: "CMP", Op1: "A", Op2: "", Length: 1, Cycles: 4, exec: cmp}, // 0xBF

	// 0xC0-0xCF
	{Name: "RNZ", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getZF() == 0) }},   // 0xC0 (11 if taken)
	{Name: "POP", Op1: "BC", Op2: "", Length: 1, Cycles: 10, exec: popOp},                                               // 0xC1
	{Name: "JNZ", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getZF() == 0) }}, // 0xC2
	{Name: "JMP", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(true) }},           // 0xC3
	{Name: "CNZ", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getZF() == 0) }}, // 0xC4 (17 if taken)
	{Name: "PUSH", Op1: "BC", Op2: "", Length: 1, Cycles: 11, exec: pushOp},                                             // 0xC5
	{Name: "ADI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: adi},                                                    // 0xC6
	{Name: "RST", Op1: "0", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                  // 0xC7
	{Name: "RZ", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getZF() == 1) }},    // 0xC8 (11 if taken)
	{Name: "RET", Op1: "", Op2: "", Length: 1, Cycles: 10, exec: func(c *CPU, _ string) { c.ret() }},                    // 0xC9
	{Name: "JZ", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getZF() == 1) }},  // 0xCA
	{Name: "JMP", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(true) }},           // 0CB *JMP
	{Name: "CZ", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getZF() == 1) }},  // 0xCC (17 if taken)
	{Name: "CALL", Op1: "", Op2: "", Length: 3, Cycles: 17, exec: func(c *CPU, _ string) { c.call() }},                  // 0xCD
	{Name: "ACI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: aci},                                                    // 0xCE
	{Name: "RST", Op1: "1", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                  // 0xCF

	// 0xD0-0xDF
	{Name: "RNC", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getCYF() == 0) }},   // 0xD0 (11 if taken)
	{Name: "POP", Op1: "DE", Op2: "", Length: 1, Cycles: 10, exec: popOp},                                                // 0xD1
	{Name: "JNC", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getCYF() == 0) }}, // 0xD2
	{Name: "OUT", Op1: "", Op2: "", Length: 2, Cycles: 10, exec: portOut},                                                // 0xD3
	{Name: "CNC", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getCYF() == 0) }}, // 0xD4 (17 if taken)
	{Name: "PUSH", Op1: "DE", Op2: "", Length: 1, Cycles: 11, exec: pushOp},                                              // 0xD5
	{Name: "SUI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: sui},                                                     // 0xD6
	{Name: "RST", Op1: "2", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                   // 0xD7
	{Name: "RC", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getCYF() == 1) }},    // 0xD8 (11 if taken)
	{Name: "RET", Op1: "", Op2: "", Length: 1, Cycles: 10, exec: func(c *CPU, _ string) { c.ret() }},                     // 0retxD9 *RET
	{Name: "JC", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getCYF() == 1) }},  // 0xDA
	{Name: "IN", Op1: "", Op2: "", Length: 2, Cycles: 10, exec: portIn},                                                  // 0xDB
	{Name: "CC", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getCYF() == 1) }},  // 0xDC (17 if taken)
	{Name: "CALL", Op1: "", Op2: "", Length: 3, Cycles: 17, exec: func(c *CPU, _ string) { c.call() }},                   // 0xDD *CALL
	{Name: "SBI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: sbi},                                                     // 0xDE
	{Name: "RST", Op1: "3", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                   // 0xDF

	// 0xE0-0xEF
	{Name: "RPO", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getPF() == 0) }},                // 0xE0 (11 if taken)
	{Name: "POP", Op1: "HL", Op2: "", Length: 1, Cycles: 10, exec: popOp},                                                            // 0xE1
	{Name: "JPO", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getPF() == 0) }},              // 0xE2
	{Name: "XTHL", Op1: "", Op2: "", Length: 1, Cycles: 18, exec: xthl},                                                              // 0xE3
	{Name: "CPO", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getPF() == 0) }},              // 0xE4 (17 if taken)
	{Name: "PUSH", Op1: "HL", Op2: "", Length: 1, Cycles: 11, exec: pushOp},                                                          // 0xE5
	{Name: "ANI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: ani},                                                                 // 0xE6
	{Name: "RST", Op1: "4", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                               // 0xE7
	{Name: "RPE", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getPF() == 1) }},                // 0xE8 (11 if taken)
	{Name: "PCHL", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.PC = uint16(c.H)<<8 | uint16(c.L) }},     // 0xE9
	{Name: "JPE", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getPF() == 1) }},              // 0xEA
	{Name: "XCHG", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.H, c.L, c.D, c.E = c.D, c.E, c.H, c.L }}, // 0xEB
	{Name: "CPE", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getPF() == 1) }},              // 0xEC (17 if taken)
	{Name: "CALL", Op1: "", Op2: "", Length: 3, Cycles: 17, exec: func(c *CPU, _ string) { c.call() }},                               // 0callxED *CALL
	{Name: "XRI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: xri},                                                                 // 0xEE
	{Name: "RST", Op1: "5", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                               // 0xEF

	// 0xF0-0xFF
	{Name: "RP", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getSF() == 0) }},             // 0xF0 (11 if taken)
	{Name: "POP", Op1: "AF", Op2: "", Length: 1, Cycles: 10, exec: popOp},                                                        // 0xF1
	{Name: "JP", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getSF() == 0) }},           // 0xF2
	{Name: "DI", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.Interrupts = false }},                  // 0xF3
	{Name: "CP", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getSF() == 0) }},           // 0xF4 (17 if taken)
	{Name: "PUSH", Op1: "AF", Op2: "", Length: 1, Cycles: 11, exec: pushOp},                                                      // 0xF5
	{Name: "ORI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: ori},                                                             // 0xF6
	{Name: "RST", Op1: "6", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                           // 0xF7
	{Name: "RM", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.retCond(c.getSF() == 1) }},             // 0xF8 (11 if taken)
	{Name: "SPHL", Op1: "", Op2: "", Length: 1, Cycles: 5, exec: func(c *CPU, _ string) { c.SP = uint16(c.H)<<8 | uint16(c.L) }}, // 0xF9
	{Name: "JM", Op1: "", Op2: "", Length: 3, Cycles: 10, exec: func(c *CPU, _ string) { c.jumpCond(c.getSF() == 1) }},           // 0xFA
	{Name: "EI", Op1: "", Op2: "", Length: 1, Cycles: 4, exec: func(c *CPU, _ string) { c.Interrupts = true }},                   // 0xFB
	{Name: "CM", Op1: "", Op2: "", Length: 3, Cycles: 11, exec: func(c *CPU, _ string) { c.callCond(c.getSF() == 1) }},           // 0xFC (17 if taken)
	{Name: "CALL", Op1: "", Op2: "", Length: 3, Cycles: 17, exec: func(c *CPU, _ string) { c.call() }},                           // 0xFD *CALL
	{Name: "CPI", Op1: "", Op2: "", Length: 2, Cycles: 7, exec: cpi},                                                             // 0xFE
	{Name: "RST", Op1: "7", Op2: "", Length: 1, Cycles: 11, exec: rst},                                                           // 0xFF
}
