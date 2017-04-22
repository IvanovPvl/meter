package main

import (
	"fmt"
	"github.com/ivanovpvl/meter/agent/df"
)

func main() {
	res, _ := df.Exec()
	fmt.Println(res)
}
