package optimizer

type CostType uint

const (
	GATE CostType = iota
	FAN_IN
)
