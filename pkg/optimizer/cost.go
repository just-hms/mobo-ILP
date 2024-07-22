package optimizer

type CostType uint

const (
	GATE_COST CostType = iota
	FAN_IN_COST
)
