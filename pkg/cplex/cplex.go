package cplex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Solve(problem string) (*CPLEXSolution, error) {

	// write the problem in a tmp input file
	input, err := os.CreateTemp("", "cplex-*.lp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(input.Name())

	_, err = input.WriteString(problem)
	if err != nil {
		return nil, err
	}
	err = input.Close()
	if err != nil {
		return nil, err
	}

	// create an output file
	output, err := os.CreateTemp("", "cplex-*.out")
	if err != nil {
		return nil, err
	}
	defer os.Remove(output.Name())

	outputName := output.Name()

	err = os.Remove(outputName)
	if err != nil {
		return nil, err
	}

	err = output.Close()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(os.Getenv("CPLEX_PATH"),
		"-c",
		fmt.Sprintf("read %s", input.Name()),
		"optimize",
		fmt.Sprintf("write %s sol", outputName),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	output, err = os.Open(outputName)
	if err != nil {
		return nil, errors.New(string(out))
	}

	solution, err := io.ReadAll(output)
	if err != nil {
		return nil, err
	}

	return extract(solution)
}

func extract(out []byte) (*CPLEXSolution, error) {
	var solution *CPLEXSolution
	err := xml.Unmarshal(out, &solution)
	if err != nil {
		return nil, err
	}
	return solution, nil
}
