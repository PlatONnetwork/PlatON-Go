package compile

import (
	"bytes"
	"encoding/binary"
	"github.com/PlatONnetwork/wagon/disasm"
	ops "github.com/PlatONnetwork/wagon/wasm/operators"
)

//func CompileWithTable(disassembly []disasm.Instr) ([]byte, []*BranchTable) {
//	buffer := new(bytes.Buffer)
//	branchTables := []*BranchTable{}
//
//	curBlockDepth := -1
//	blocks := make(map[int]*block) // maps nesting depths (labels) to blocks
//
//	blocks[-1] = &block{}
//
//	for _, instr := range disassembly {
//		if instr.Unreachable {
//			continue
//		}
//		switch instr.Op.Code {
//		//case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
//		case ops.I32Load, ops.I64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
//			// memory_immediate has two fields, the alignment and the offset.
//			// The former is simply an optimization hint and can be safely
//			// discarded.
//			instr.Immediates = []interface{}{instr.Immediates[1].(uint32)}
//		case ops.If:
//			curBlockDepth++
//			buffer.WriteByte(OpJmpZ)
//			blocks[curBlockDepth] = &block{
//				ifBlock:        true,
//				elseAddrOffset: int64(buffer.Len()),
//			}
//			// the address to jump to if the condition for `if` is false
//			// (i.e when the value on the top of the stack is 0)
//			binary.Write(buffer, binary.LittleEndian, int64(0))
//			continue
//		case ops.Loop:
//			// there is no condition for entering a loop block
//			curBlockDepth++
//			blocks[curBlockDepth] = &block{
//				offset:    int64(buffer.Len()),
//				ifBlock:   false,
//				loopBlock: true,
//				discard:   *instr.NewStack,
//			}
//			continue
//		case ops.Block:
//			curBlockDepth++
//			blocks[curBlockDepth] = &block{
//				ifBlock: false,
//				discard: *instr.NewStack,
//			}
//			continue
//		case ops.Else:
//			ifInstr := disassembly[instr.Block.ElseIfIndex] // the corresponding `if` instruction for this else
//			if ifInstr.NewStack != nil && ifInstr.NewStack.StackTopDiff != 0 {
//				// add code for jumping out of a taken if branch
//				if ifInstr.NewStack.PreserveTop {
//					buffer.WriteByte(OpDiscardPreserveTop)
//				} else {
//					buffer.WriteByte(OpDiscard)
//				}
//				binary.Write(buffer, binary.LittleEndian, ifInstr.NewStack.StackTopDiff)
//			}
//			buffer.WriteByte(OpJmp)
//			ifBlockEndOffset := int64(buffer.Len())
//			binary.Write(buffer, binary.LittleEndian, int64(0))
//
//			curOffset := int64(buffer.Len())
//			ifBlock := blocks[curBlockDepth]
//			code := buffer.Bytes()
//
//			buffer = patchOffsetNoTarget(code, ifBlock.elseAddrOffset, curOffset)
//			// this is no longer an if block
//			ifBlock.ifBlock = false
//			ifBlock.patchOffsets = append(ifBlock.patchOffsets, ifBlockEndOffset)
//			continue
//		case ops.End:
//			depth := curBlockDepth
//			block := blocks[depth]
//
//			if instr.NewStack.StackTopDiff != 0 {
//				// when exiting a block, discard elements to
//				// restore stack height.
//				if instr.NewStack.PreserveTop {
//					// this is true when the block has a
//					// signature, and therefore pushes
//					// a value on to the stack
//					buffer.WriteByte(OpDiscardPreserveTop)
//				} else {
//					buffer.WriteByte(OpDiscard)
//				}
//				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)
//			}
//
//			if !block.loopBlock { // is a normal block
//				block.offset = int64(buffer.Len())
//				if block.ifBlock {
//					code := buffer.Bytes()
//					buffer = patchOffsetNoTarget(code, block.elseAddrOffset, int64(block.offset))
//				}
//			}
//
//			for _, offset := range block.patchOffsets {
//				code := buffer.Bytes()
//				buffer = patchOffsetNoTarget(code, offset, block.offset)
//			}
//
//			for _, table := range block.branchTables {
//				table.patchTableNoTarget(table.blocksLen-depth-1, int64(block.offset))
//			}
//
//			delete(blocks, curBlockDepth)
//			curBlockDepth--
//			continue
//		case ops.Br:
//			if instr.NewStack != nil && instr.NewStack.StackTopDiff != 0 {
//				if instr.NewStack.PreserveTop {
//					buffer.WriteByte(OpDiscardPreserveTop)
//				} else {
//					buffer.WriteByte(OpDiscard)
//				}
//				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)
//			}
//			buffer.WriteByte(OpJmp)
//			label := int(instr.Immediates[0].(uint32))
//			block := blocks[curBlockDepth-int(label)]
//			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
//			// write the jump address
//			binary.Write(buffer, binary.LittleEndian, int64(0))
//			continue
//		case ops.BrIf:
//			buffer.WriteByte(OpJmpNz)
//			label := int(instr.Immediates[0].(uint32))
//			block := blocks[curBlockDepth-int(label)]
//			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
//			// write the jump address
//			binary.Write(buffer, binary.LittleEndian, int64(0))
//
//			var stackTopDiff int64
//			// write whether we need to preserve the top
//			if instr.NewStack == nil || !instr.NewStack.PreserveTop || instr.NewStack.StackTopDiff == 0 {
//				buffer.WriteByte(byte(0))
//			} else {
//				stackTopDiff = instr.NewStack.StackTopDiff
//				buffer.WriteByte(byte(1))
//			}
//			// write the number of elements on the stack we need to discard
//			binary.Write(buffer, binary.LittleEndian, stackTopDiff)
//			continue
//		case ops.BrTable:
//			branchTable := &BranchTable{
//				// we subtract one for the implicit block created by
//				// the function body
//				blocksLen: len(blocks) - 1,
//			}
//			targetCount := instr.Immediates[0].(uint32)
//			branchTable.Targets = make([]Target, targetCount)
//			for i := range branchTable.Targets {
//				// The first immediates is the number of targets, so we ignore that
//				label := int64(instr.Immediates[i+1].(uint32))
//				branchTable.Targets[i].Addr = label
//				branch := instr.Branches[i]
//
//				branchTable.Targets[i].Return = branch.IsReturn
//				branchTable.Targets[i].Discard = branch.StackTopDiff
//				branchTable.Targets[i].PreserveTop = branch.PreserveTop
//			}
//			defaultLabel := int64(instr.Immediates[len(instr.Immediates)-1].(uint32))
//			branchTable.DefaultTarget.Addr = defaultLabel
//			defaultBranch := instr.Branches[targetCount]
//			branchTable.DefaultTarget.Return = defaultBranch.IsReturn
//			branchTable.DefaultTarget.Discard = defaultBranch.StackTopDiff
//			branchTable.DefaultTarget.PreserveTop = defaultBranch.PreserveTop
//			branchTables = append(branchTables, branchTable)
//			for _, block := range blocks {
//				block.branchTables = append(block.branchTables, branchTable)
//			}
//
//			buffer.WriteByte(ops.BrTable)
//			binary.Write(buffer, binary.LittleEndian, int64(len(branchTables)-1))
//		}
//
//		buffer.WriteByte(instr.Op.Code)
//		for _, imm := range instr.Immediates {
//			err := binary.Write(buffer, binary.LittleEndian, imm)
//			if err != nil {
//				panic(err)
//			}
//		}
//	}
//
//	// writing nop as the last instructions allows us to branch out of the
//	// function (ie, return)
//	addr := buffer.Len()
//	buffer.WriteByte(ops.Nop)
//
//	// patch all references to the "root" block of the function body
//	for _, offset := range blocks[-1].patchOffsets {
//		code := buffer.Bytes()
//		buffer = patchOffsetNoTarget(code, offset, int64(addr))
//	}
//
//	for _, table := range branchTables {
//		table.patchedAddrs = nil
//	}
//	return buffer.Bytes(), branchTables
//}
//
//
//// replace the address starting at start with addr
//func patchOffsetNoTarget(code []byte, start int64, addr int64) *bytes.Buffer {
//	var shift uint
//	for i := int64(0); i < 8; i++ {
//		code[start+i] = byte(addr >> shift)
//		shift += 8
//	}
//
//	buf := new(bytes.Buffer)
//	buf.Write(code)
//	return buf
//}
//
//func (table *BranchTable) patchTableNoTarget(block int, addr int64) {
//	if block < 0 {
//		panic("Invalid block value")
//	}
//
//	for i, target := range table.Targets {
//		if !table.isAddr(target.Addr) && target.Addr == int64(block) {
//			table.Targets[i].Addr = addr
//		}
//	}
//
//	if table.DefaultTarget.Addr == int64(block) {
//		table.DefaultTarget.Addr = addr
//	}
//	table.patchedAddrs = append(table.patchedAddrs, addr)
//}


// Compile rewrites WebAssembly bytecode from its disassembly.
// TODO(vibhavp): Add options for optimizing code. Operators like i32.reinterpret/f32
// are no-ops, and can be safely removed.
func CompileWithTable(disassembly []disasm.Instr) ([]byte, []*BranchTable) {
	buffer := new(bytes.Buffer)
	branchTables := []*BranchTable{}

	curBlockDepth := -1
	blocks := make(map[int]*block) // maps nesting depths (labels) to blocks

	blocks[-1] = &block{}

	for _, instr := range disassembly {
		if instr.Unreachable {
			continue
		}
		switch instr.Op.Code {
		//case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
		case ops.I32Load, ops.I64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			// memory_immediate has two fields, the alignment and the offset.
			// The former is simply an optimization hint and can be safely
			// discarded.
			instr.Immediates = []interface{}{instr.Immediates[1].(uint32)}
		case ops.If:
			curBlockDepth++
			buffer.WriteByte(OpJmpZ)
			blocks[curBlockDepth] = &block{
				ifBlock:        true,
				elseAddrOffset: int64(buffer.Len()),
			}
			// the address to jump to if the condition for `if` is false
			// (i.e when the value on the top of the stack is 0)
			binary.Write(buffer, binary.LittleEndian, int64(0))
			continue
		case ops.Loop:
			// there is no condition for entering a loop block
			curBlockDepth++
			blocks[curBlockDepth] = &block{
				offset:    int64(buffer.Len()),
				ifBlock:   false,
				loopBlock: true,
				discard:   *instr.NewStack,
			}
			continue
		case ops.Block:
			curBlockDepth++
			blocks[curBlockDepth] = &block{
				ifBlock: false,
				discard: *instr.NewStack,
			}
			continue
		case ops.Else:
			ifInstr := disassembly[instr.Block.ElseIfIndex] // the corresponding `if` instruction for this else
			if ifInstr.NewStack != nil && ifInstr.NewStack.StackTopDiff != 0 {
				// add code for jumping out of a taken if branch
				if ifInstr.NewStack.PreserveTop {
					buffer.WriteByte(OpDiscardPreserveTop)
				} else {
					buffer.WriteByte(OpDiscard)
				}
				binary.Write(buffer, binary.LittleEndian, ifInstr.NewStack.StackTopDiff)
			}
			buffer.WriteByte(OpJmp)
			ifBlockEndOffset := int64(buffer.Len())
			binary.Write(buffer, binary.LittleEndian, int64(0))

			curOffset := int64(buffer.Len())
			ifBlock := blocks[curBlockDepth]
			code := buffer.Bytes()

			buffer = patchOffsetNoTarget(code, ifBlock.elseAddrOffset, curOffset)
			// this is no longer an if block
			ifBlock.ifBlock = false
			ifBlock.patchOffsets = append(ifBlock.patchOffsets, ifBlockEndOffset)
			continue
		case ops.End:
			depth := curBlockDepth
			block := blocks[depth]

			if instr.NewStack.StackTopDiff != 0 {
				// when exiting a block, discard elements to
				// restore stack height.
				if instr.NewStack.PreserveTop {
					// this is true when the block has a
					// signature, and therefore pushes
					// a value on to the stack
					buffer.WriteByte(OpDiscardPreserveTop)
				} else {
					buffer.WriteByte(OpDiscard)
				}
				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)
			}

			if !block.loopBlock { // is a normal block
				block.offset = int64(buffer.Len())
				if block.ifBlock {
					code := buffer.Bytes()
					buffer = patchOffsetNoTarget(code, block.elseAddrOffset, int64(block.offset))
				}
			}

			for _, offset := range block.patchOffsets {
				code := buffer.Bytes()
				buffer = patchOffsetNoTarget(code, offset, block.offset)
			}

			for _, table := range block.branchTables {
				table.patchTableNoTarget(table.blocksLen-depth-1, int64(block.offset))
			}

			delete(blocks, curBlockDepth)
			curBlockDepth--
			continue
		case ops.Br:
			if instr.NewStack != nil && instr.NewStack.StackTopDiff != 0 {
				if instr.NewStack.PreserveTop {
					buffer.WriteByte(OpDiscardPreserveTop)
				} else {
					buffer.WriteByte(OpDiscard)
				}
				binary.Write(buffer, binary.LittleEndian, instr.NewStack.StackTopDiff)
			}
			buffer.WriteByte(OpJmp)
			label := int(instr.Immediates[0].(uint32))
			block := blocks[curBlockDepth-int(label)]
			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
			// write the jump address
			binary.Write(buffer, binary.LittleEndian, int64(0))
			continue
		case ops.BrIf:
			buffer.WriteByte(OpJmpNz)
			label := int(instr.Immediates[0].(uint32))
			block := blocks[curBlockDepth-int(label)]
			block.patchOffsets = append(block.patchOffsets, int64(buffer.Len()))
			// write the jump address
			binary.Write(buffer, binary.LittleEndian, int64(0))

			var stackTopDiff int64
			// write whether we need to preserve the top
			if instr.NewStack == nil || !instr.NewStack.PreserveTop || instr.NewStack.StackTopDiff == 0 {
				buffer.WriteByte(byte(0))
			} else {
				stackTopDiff = instr.NewStack.StackTopDiff
				buffer.WriteByte(byte(1))
			}
			// write the number of elements on the stack we need to discard
			binary.Write(buffer, binary.LittleEndian, stackTopDiff)
			continue
		case ops.BrTable:
			branchTable := &BranchTable{
				// we subtract one for the implicit block created by
				// the function body
				blocksLen: len(blocks) - 1,
			}
			targetCount := instr.Immediates[0].(uint32)
			branchTable.Targets = make([]Target, targetCount)
			for i := range branchTable.Targets {
				// The first immediates is the number of targets, so we ignore that
				label := int64(instr.Immediates[i+1].(uint32))
				branchTable.Targets[i].Addr = label
				branch := instr.Branches[i]

				branchTable.Targets[i].Return = branch.IsReturn
				branchTable.Targets[i].Discard = branch.StackTopDiff
				branchTable.Targets[i].PreserveTop = branch.PreserveTop
			}
			defaultLabel := int64(instr.Immediates[len(instr.Immediates)-1].(uint32))
			branchTable.DefaultTarget.Addr = defaultLabel
			defaultBranch := instr.Branches[targetCount]
			branchTable.DefaultTarget.Return = defaultBranch.IsReturn
			branchTable.DefaultTarget.Discard = defaultBranch.StackTopDiff
			branchTable.DefaultTarget.PreserveTop = defaultBranch.PreserveTop
			branchTables = append(branchTables, branchTable)
			for _, block := range blocks {
				block.branchTables = append(block.branchTables, branchTable)
			}

			buffer.WriteByte(ops.BrTable)
			binary.Write(buffer, binary.LittleEndian, int64(len(branchTables)-1))
		}

		buffer.WriteByte(instr.Op.Code)
		for _, imm := range instr.Immediates {
			err := binary.Write(buffer, binary.LittleEndian, imm)
			if err != nil {
				panic(err)
			}
		}
	}

	// writing nop as the last instructions allows us to branch out of the
	// function (ie, return)
	addr := buffer.Len()
	buffer.WriteByte(ops.Nop)

	// patch all references to the "root" block of the function body
	for _, offset := range blocks[-1].patchOffsets {
		code := buffer.Bytes()
		buffer = patchOffsetNoTarget(code, offset, int64(addr))
	}

	for _, table := range branchTables {
		table.patchedAddrs = nil
	}
	return buffer.Bytes(), branchTables
}

// replace the address starting at start with addr
func patchOffsetNoTarget(code []byte, start int64, addr int64) *bytes.Buffer {
	var shift uint
	for i := int64(0); i < 8; i++ {
		code[start+i] = byte(addr >> shift)
		shift += 8
	}

	buf := new(bytes.Buffer)
	buf.Write(code)
	return buf
}

func (table *BranchTable) patchTableNoTarget(block int, addr int64) {
	if block < 0 {
		panic("Invalid block value")
	}

	for i, target := range table.Targets {
		if !table.isAddr(target.Addr) && target.Addr == int64(block) {
			table.Targets[i].Addr = addr
		}
	}

	if table.DefaultTarget.Addr == int64(block) {
		table.DefaultTarget.Addr = addr
	}
	table.patchedAddrs = append(table.patchedAddrs, addr)
}