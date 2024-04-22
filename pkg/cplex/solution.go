package cplex

import "encoding/xml"

type CPLEXSolution struct {
	XMLName           xml.Name     `xml:"CPLEXSolution"`
	Header            Header       `xml:"header"`
	Quality           Quality      `xml:"quality"`
	LinearConstraints []Constraint `xml:"linearConstraints>constraint"`
	Variables         []Variable   `xml:"variables>variable"`
}

type Header struct {
	ProblemName          string  `xml:"problemName,attr"`
	SolutionName         string  `xml:"solutionName,attr"`
	ObjectiveValue       float64 `xml:"objectiveValue,attr"`
	SolutionStatusString string  `xml:"solutionStatusString,attr"`
}

type Quality struct {
	MaxIntInfeas    float64 `xml:"maxIntInfeas,attr"`
	MaxPrimalInfeas float64 `xml:"maxPrimalInfeas,attr"`
}

type Constraint struct {
	Name  string  `xml:"name,attr"`
	Slack float64 `xml:"slack,attr"`
}

type Variable struct {
	Name  string  `xml:"name,attr"`
	Value float64 `xml:"value,attr"`
}
