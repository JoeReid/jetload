package specfile

import (
	"fmt"
	"os"
	"strconv"
	"text/template"
)

var basicFuncs = template.FuncMap{
	"seq": seq,
	"env": env,
	"mod": mod,
	"int": asInt,
}

func seq(n int) []int {
	seq := make([]int, n)
	for i := range seq {
		seq[i] = i
	}
	return seq
}

func env(key string) string {
	return os.Getenv(key)
}

func mod(b, a int) int {
	return a % b
}

func asInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case string:
		return strconv.Atoi(v)

	default:
		return 0, fmt.Errorf("invalid type %T", v)
	}
}
