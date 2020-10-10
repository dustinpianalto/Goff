package fun

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/rpnparse"
)

var InterleaveCommand = &disgoman.Command{
	Name:                "encode",
	Aliases:             []string{"e"},
	Description:         "Encode 2 numbers",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              interleaveFunc,
}

func interleaveFunc(ctx disgoman.Context, args []string) {
	if len(args) == 2 {
		x, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return
		}
		y, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return
		}
		var z = int64(0)
		for i := 0; i < 64; i++ {
			x_masked_i := x & (1 << i)
			y_masked_i := y & (1 << i)

			z |= x_masked_i << i
			z |= y_masked_i << (i + 1)
		}
		ctx.Send(fmt.Sprintf("%v", z))
	}
}

var DeinterleaveCommand = &disgoman.Command{
	Name:                "decode",
	Aliases:             []string{"d"},
	Description:         "Decode 1 number into 2",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              deinterleaveFunc,
}

func deinterleaveFunc(ctx disgoman.Context, args []string) {
	if len(args) == 1 {
		z, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return
		}
		var x = int64(0)
		var y = int64(0)
		i := 0
		for z > 0 {
			x |= (z & 1) << i
			z >>= 1
			y |= (z & 1) << i
			z >>= 1
			i++
		}
		ctx.Send(fmt.Sprintf("(%v, %v)", x, y))
	}
}

var GenerateRPNCommand = &disgoman.Command{
	Name:                "RPN",
	Aliases:             []string{"rpn"},
	Description:         "Convert infix to rpn",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              generateRPNFunc,
}

func generateRPNFunc(ctx disgoman.Context, args []string) {
	rpn, err := rpnparse.GenerateRPN(args)
	if err != nil {
		ctx.Send(err.Error())
		return
	}
	ctx.Send(rpn)
}

var ParseRPNCommand = &disgoman.Command{
	Name:                "ParseRPN",
	Aliases:             []string{"PRPN", "prpn"},
	Description:         "Parse RPN string and return the result",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              parseRPNFunc,
}

func parseRPNFunc(ctx disgoman.Context, args []string) {
	res, err := rpnparse.ParseRPN(args)
	if err != nil {
		ctx.Send(err.Error())
		return
	}
	ctx.Send(fmt.Sprintf("The result is: %v", res))
}

var SolveCommand = &disgoman.Command{
	Name:                "solve",
	Aliases:             []string{"math", "infix"},
	Description:         "Solve infix equation and return the result",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              solveFunc,
}

func solveFunc(ctx disgoman.Context, args []string) {
	rpn, err := rpnparse.GenerateRPN(args)
	if err != nil {
		ctx.Send(err.Error())
		return
	}
	res, err := rpnparse.ParseRPN(strings.Split(rpn, " "))
	if err != nil {
		ctx.Send(err.Error())
		return
	}
	ctx.Send(fmt.Sprintf("The result is: %v", res))
}
