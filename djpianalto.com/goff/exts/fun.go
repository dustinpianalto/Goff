package exts

import (
	"fmt"
	"github.com/dustinpianalto/disgoman"
	"strconv"
)

func interleave(ctx disgoman.Context, args []string) {
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
		ctx.Send(string(z))
	}
}

func deinterleave(ctx disgoman.Context, args []string) {
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
