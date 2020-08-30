package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type FStack []float64

func (s *FStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *FStack) Push(op float64) {
	*s = append(*s, op)
}

func (s *FStack) Pop() (float64, bool) {
	if s.IsEmpty() {
		return 0, false
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}

func (s *FStack) PopTwo() (float64, float64, bool) {
	if s.IsEmpty() || len(*s) < 2 {
		return 0, 0, false
	}
	index := len(*s) - 1
	b := (*s)[index]
	a := (*s)[index-1]
	*s = (*s)[:index-1]
	return a, b, true

}

func (s *FStack) Top() float64 {
	if s.IsEmpty() {
		return 0
	}
	return (*s)[len(*s)-1]
}

func ParseRPN(args []string) (float64, error) {
	s := FStack{}
	for _, token := range args {
		switch token {
		case "+":
			if a, b, ok := s.PopTwo(); ok {
				s.Push(a + b)
			} else {
				return 0, fmt.Errorf("not enough operands on stack for +: %v", s)
			}
		case "-":
			if a, b, ok := s.PopTwo(); ok {
				s.Push(a - b)
			} else {
				return 0, fmt.Errorf("not enough operands on stack for -: %v", s)
			}
		case "*":
			if a, b, ok := s.PopTwo(); ok {
				s.Push(a * b)
			} else {
				return 0, fmt.Errorf("not enough operands on stack for *: %v", s)
			}
		case "/":
			if a, b, ok := s.PopTwo(); ok {
				s.Push(a / b)
			} else {
				return 0, fmt.Errorf("not enough operands on stack for /: %v", s)
			}
		case "%":
			if a, b, ok := s.PopTwo(); ok {
				s.Push(math.Mod(a, b))
			} else {
				return 0, fmt.Errorf("not enough operands on stack for %: %v", s)
			}
		default:
			f, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			s.Push(f)
		}
	}
	if res, ok := s.Pop(); ok {
		return res, nil
	}
	return 0, errors.New("no result")
}
