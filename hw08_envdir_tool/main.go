package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		err error
		env Environment
		out string
	)
	path := os.Args[1]
	cmd := os.Args[2:]

	env, err = ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if out, err = RunCmd(cmd, env); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out)
}
