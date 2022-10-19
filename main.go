/*
Copyright © 2021 cloud.native
*/
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/kakao/detek/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "unhandled panic occured: %v", r)
			fmt.Fprintf(os.Stderr, "stacktrace from panic: \n %v", string(debug.Stack()))
			os.Exit(-1)
		}
	}()
	cmd.Execute()
}
