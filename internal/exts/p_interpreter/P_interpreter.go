package p_interpreter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dustinpianalto/disgoman"
)

var PCommand = &disgoman.Command{
	Name:                "P",
	Aliases:             nil,
	Description:         "Interpret a P\" program and return the results",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              pCommandFunc,
}

func pCommandFunc(ctx disgoman.Context, args []string) {
	input := strings.Join(args, "")
	const LENGTH = 1999
	var mem [LENGTH]byte
	pointer := 0
	l := 0
	for i := 0; i < len(input); i++ {
		if input[i] == 'L' {
			if pointer == 0 {
				pointer = LENGTH - 1
			} else {
				pointer--
			}
		} else if input[i] == 'R' {
			if pointer == LENGTH-1 {
				pointer = 0
			} else {
				pointer++
			}
		} else if input[i] == '+' {
			mem[pointer]++
		} else if input[i] == '-' {
			mem[pointer]--
		} else if input[i] == '(' {
			if mem[pointer] == 0 {
				i++
				for l > 0 || input[i] != ')' {
					if input[i] == '(' {
						l++
					}
					if input[i] == ')' {
						l--
					}
					i++
				}
			}
		} else if input[i] == ')' {
			if mem[pointer] != 0 {
				i--
				for l > 0 || input[i] != '(' {
					if input[i] == ')' {
						l++
					}
					if input[i] == '(' {
						l--
					}
					i--
				}
			}
		} else {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: fmt.Sprintf("Invalid Character: %v", input[i]),
				Error:   errors.New("invalid character"),
			}
			return
		}
	}
	var out []byte
	for _, i := range mem {
		if i != 0 {
			out = append(out, i)
		}
	}
	_, err := ctx.Send(string(out))
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Couldn't send results",
			Error:   err,
		}
		return
	}
	return
}
