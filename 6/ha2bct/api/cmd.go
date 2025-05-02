package main

import (
	"flag"
	"fmt"
	"ha2bct/usecase"
	"os"
)

func main() {
	translateFile()
}

func translateFile() {
	translateCmd := flag.NewFlagSet("translate", flag.ExitOnError)
	path := translateCmd.String("p", "missed", "path")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "translate":
		if err := translateCmd.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		outFilename, err := usecase.NewHackAssemblerTo16BitFileTranslator().TranslateAll(*path)
		if err != nil {
			panic(err)
		}
		fmt.Println("Result file:" + outFilename)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd.go <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  translate -p src/filename -> translate .asm file to .hack file (assembler to binary code)")
}
