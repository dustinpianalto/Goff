package exts

import (
	"errors"
	"fmt"
	"github.com/dustinpianalto/disgoman"
	"strings"
)

func pCommand(ctx disgoman.Context, args []string) error {
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
			ctx.Send(fmt.Sprintf("Invalid Character: %v", input[i]))
			return errors.New("invalid character")
		}
	}
	var out []byte
	for _, i := range mem {
		if i != 0 {
			out = append(out, i)
		}
	}
	fmt.Println(out)
	_, err := ctx.Send(string(out))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
