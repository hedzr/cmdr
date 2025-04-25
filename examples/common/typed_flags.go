package common

import (
	"time"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddTypedFlags(c cli.CommandBuilder) {
	c.Flg("bool", "b").
		Default(false).
		Description("A bool flag", "").
		Group("0999.Bool").
		Build()

	c.Flg("int", "i").
		Default(1).
		Description("A int flag", "").
		Group("1000.Integer").
		Build()

	c.Flg("int64", "i64").
		Default(int64(2)).
		Description("A int64 flag", "").
		Group("1000.Integer").
		Build()

	c.Flg("uint", "u").
		Default(uint(3)).
		Description("A uint flag", "").
		Group("1000.Integer").
		Build()

	c.Flg("uint64", "u64").
		Default(uint64(4)).
		Description("A uint64 flag", "").
		Group("1000.Integer").
		Build()

	c.Flg("float32", "f", "float").
		Default(float32(2.71828)).
		Description("A float32 flag with 'e' value", "").
		Group("2000.Float").
		Build()

	c.Flg("float64", "f64").
		Default(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
		Description("A float64 flag with a `PI` value", "").
		Group("2000.Float").
		Build()

	c.Flg("complex64", "c64").
		Default(complex64(3.14+9i)).
		Description("A complex64 flag", "").
		Group("2010.Complex").
		Build()

	c.Flg("complex128", "c128").
		Default(complex128(3.14+9i)).
		Description("A complex128 flag", "").
		Group("2010.Complex").
		Build()

	// a set of booleans

	c.Flg("single", "sng").
		Default(false).
		Description("A bool flag: single", "").
		Group("Booleans").
		EnvVars("").
		Build()

	c.Flg("double", "dbl").
		Default(false).
		Description("A bool flag: double", "").
		Group("Booleans").
		EnvVars("").
		Build()

	c.Flg("norway", "nw").
		Default(false).
		Description("A bool flag: norway", "").
		Group("Booleans").
		Build()

	c.Flg("mongo", "mongo").
		Default(false).
		Description("A bool flag: mongo", "").
		Group("Booleans").
		Build()

	// durations

	c.Flg("duration", "dur").
		Default(3*time.Second).
		Description("A duration flag", "").
		Group("Time & Duration").
		Build()

	// With(func(b cli.CommandBuilder) {
	// 	b.Flg("duration", "dur").
	// 		Default(time.Second * 5).
	// 		Description("a duration var").
	// 		Build()
	// })
}
