package main

/* instructions populate 'ins' map. Key is the instruction code and the value
 * is the function to call according to it's registor type
 */
func instructions() (ins map[uint32]interface{}) {
	ins = make(map[uint32]interface{})

	ins[0x00] = TypeU{fn: "Add"}
	ins[0x02] = TypeU{fn: "Sub"}
	ins[0x04] = TypeU{fn: "Mul"}
	ins[0x06] = TypeU{fn: "Div"}
	ins[0x0A] = TypeU{fn: "Shl"}
	ins[0x0B] = TypeU{fn: "Shr"}

	ins[0x0C] = TypeU{fn: "And"}
	ins[0x0E] = TypeU{fn: "Not"}
	ins[0x10] = TypeU{fn: "Or"}
	ins[0x12] = TypeU{fn: "Xor"}

	ins[0x01] = TypeF{fn: "Addi"}
	ins[0x03] = TypeF{fn: "Subi"}
	ins[0x05] = TypeF{fn: "Muli"}
	ins[0x07] = TypeF{fn: "Divi"}
	ins[0x08] = TypeF{fn: "Cmp"}
	ins[0x09] = TypeF{fn: "Cmpi"}
	ins[0x0D] = TypeF{fn: "Andi"}
	ins[0x0F] = TypeF{fn: "Noti"}
	ins[0x11] = TypeF{fn: "Ori"}
	ins[0x13] = TypeF{fn: "Xori"}

	ins[0x14] = TypeF{fn: "Ldw"}
	ins[0x15] = TypeF{fn: "Ldb"}
	ins[0x16] = TypeF{fn: "Stw"}
	ins[0x17] = TypeF{fn: "Stb"}

	ins[0x18] = TypeF{fn: "Lfr"}
	ins[0x19] = TypeF{fn: "Sfr"}
	ins[0x21] = TypeF{fn: "Call"}
	ins[0x22] = TypeF{fn: "Isr"}
	ins[0x23] = TypeF{fn: "Ret"}

	ins[0x1A] = TypeS{fn: "Bun"}
	ins[0x1B] = TypeS{fn: "Beq"}
	ins[0x1C] = TypeS{fn: "Blt"}
	ins[0x1D] = TypeS{fn: "Bgt"}
	ins[0x1E] = TypeS{fn: "Bne"}
	ins[0x1F] = TypeS{fn: "Ble"}
	ins[0x20] = TypeS{fn: "Bge"}
	ins[0x3F] = TypeS{fn: "Int"}
	return
}
