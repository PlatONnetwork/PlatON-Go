package compiler

// Gas costs
const (
	GasQuickStep	uint64 = 2
	GasFastestSetp	uint64 = 3

	// ...
)


func (c *SSAFunctionCompiler) InsertGasCounters(gp GasPolicy) {
	cfg := c.NewCFGraph()

	for i, _ := range cfg.Blocks {
		blk := &cfg.Blocks[i]
		totalCost := int64(0)
		for _, ins := range blk.Code {
			totalCost += gp.GetCost(ins.Op)
			if totalCost < 0 {
				panic("total cost overflow")
			}
		}

		if totalCost != 0 {
			blk.Code = append([]Instr{
				buildInstr(0, "add_gas", []int64{totalCost}, []TyValueID{}),
			}, blk.Code...)
		}
	}
	c.Code = cfg.ToInsSeq()
}
