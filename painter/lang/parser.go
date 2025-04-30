package lang

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct{}

func parseFloatArgs(args []string) ([]float64, bool) {
	var res []float64
	for _, s := range args {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, false
		}
		res = append(res, f)
	}
	return res, true
}

func (p *Parser) Parse(r io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(r)
	var ops []painter.Operation

	loop := &painter.Loop{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		cmd := tokens[0]
		args := tokens[1:]

		switch cmd {
		case "white":
			ops = append(ops, painter.NewWhiteOp(loop))
		case "green":
			ops = append(ops, painter.NewGreenOp(loop))
		case "update":
			ops = append(ops, painter.NewUpdateOp(loop))
		case "reset":
			ops = append(ops, painter.NewResetOp(loop))
		case "bgrect":
			if len(args) != 4 {
				continue
			}
			if nums, ok := parseFloatArgs(args); ok {
				ops = append(ops, painter.NewBgRectOp(loop, nums[0], nums[1], nums[2], nums[3]))
			}
		case "figure":
			if len(args) != 2 {
				continue
			}
			if nums, ok := parseFloatArgs(args); ok {
				ops = append(ops, painter.NewFigureOp(loop, nums[0], nums[1]))
			}
		case "move":
			if len(args) != 2 {
				continue
			}
			if nums, ok := parseFloatArgs(args); ok {
				ops = append(ops, painter.NewMoveOp(loop, nums[0], nums[1]))
			}
		}
	}

	return ops, nil
}
