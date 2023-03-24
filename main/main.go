package main

import (
	"lookup"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		lookup.Run(".", 2)
	} else if len(os.Args) == 2 {
		lookup.Run(os.Args[1], 2)
	} else if len(os.Args) == 3 {
		depth, _ := strconv.Atoi(os.Args[2])
		lookup.Run(os.Args[1], depth)
	} else {
		panic("Wrong input args")
	}
}

//.\lookup\main\main.exe
