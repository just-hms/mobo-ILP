package mobo

import (
	"strings"

	"github.com/just-hms/mobo/pkg/qm/cube"
)

type Circuit []*cube.Cube

func (c Circuit) Display(size uint) string {
	s := []string{}
	for _, gate := range c {
		s = append(s, gate.Display(size))
	}

	return strings.Join(s, " + ")
}
