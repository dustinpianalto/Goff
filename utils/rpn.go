package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Operator struct {
	Token       string
	Precedence  int
	Association string
}

func (o Operator) HasHigherPrecedence(t Operator) bool {
	return o.Precedence < t.Precedence // lower number is higher precedence
}

func (o Operator) HasEqualPrecedence(t Operator) bool {
	return o.Precedence == t.Precedence
}

func (o Operator) IsLeftAssociative() bool {
	return o.Association == "left"
}

var operators = map[string]Operator{
	"+": Operator{
		Token:       "+",
		Precedence:  4,
		Association: "left",
	},
	"-": Operator{
		Token:       "-",
		Precedence:  4,
		Association: "left",
	},
	"*": Operator{
		Token:       "*",
		Precedence:  3,
		Association: "left",
	},
	"/": Operator{
		Token:       "/",
		Precedence:  3,
		Association: "left",
	},
	"%": Operator{
		Token:       "%",
		Precedence:  3,
		Association: "left",
	},
	"(": Operator{
		Token:       "(",
		Precedence:  1,
		Association: "left",
	},
	")": Operator{
		Token:       ")",
		Precedence:  1,
		Association: "left",
	},
}

type Stack []Operator

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(op Operator) {
	*s = append(*s, op)
}

func (s *Stack) Pop() (Operator, bool) {
	if s.IsEmpty() {
		return Operator{}, false
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}

func (s *Stack) Top() Operator {
	if s.IsEmpty() {
		return Operator{}
	}
	return (*s)[len(*s)-1]
}

func GenerateRPN(tokens []string) (string, error) {
	output := ""
	s := Stack{}
	for _, token := range tokens {
		err := processToken(token, &s, &output)
		if err != nil {
			return "", err
		}
	}
	for !s.IsEmpty() {
		ele, _ := s.Pop()
		output += " " + ele.Token
	}

	return strings.TrimSpace(output), nil
}

func processToken(t string, s *Stack, o *string) error {
	if _, err := strconv.Atoi(t); err == nil {
		*o += " " + t
		return nil
	} else if op, ok := operators[t]; ok {
		if op.Token == "(" {
			s.Push(op)
		} else if op.Token == ")" {
			if s.IsEmpty() {
				return fmt.Errorf("mismatched parentheses")
			}
			for s.Top().Token != "(" {
				if ele, ok := s.Pop(); ok {
					*o += " " + ele.Token
				} else {
					return fmt.Errorf("mismatched parentheses")
				}
				if s.IsEmpty() {
					break
				}
			}
			s.Pop() // Pop and discard the (
		} else if !s.IsEmpty() {
			for {
				if (s.Top().HasHigherPrecedence(op) ||
					(s.Top().HasEqualPrecedence(op) &&
						op.IsLeftAssociative())) &&
					s.Top().Token != "(" {
					if ele, ok := s.Pop(); ok {
						*o += " " + ele.Token
						if s.IsEmpty() {
							break
						}
						continue
					} else {
						break
					}
				}
				break
			}
			s.Push(op)
		} else {
			s.Push(op)
		}
		return nil
	}
	return fmt.Errorf("invalid character %s", t)
}
