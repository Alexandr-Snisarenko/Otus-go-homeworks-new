package main

import (
	"flag"
	"fmt"
)

var (
	vFrom, vTo      string
	vLimit, vOffset int64
	vRewrite        bool
)

func init() {
	flag.StringVar(&vFrom, "from", "", "file to read from")
	flag.StringVar(&vTo, "to", "", "file to write to")
	flag.Int64Var(&vLimit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&vOffset, "offset", 0, "offset in input file")
	flag.BoolVar(&vRewrite, "rewrite", false, "rewrite file if exists ")
}

func main() {
	flag.Parse()

	if err := Copy(vFrom, vTo, vOffset, vLimit, vRewrite); err != nil {
		fmt.Println("Error:  ", err)
	}
}
