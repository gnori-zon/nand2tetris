package main

import (
	"fmt"
	"ha2bct/core"
)

func main() {
	translator := core.New16bitTranslator()
	translatedRows, err := translator.Translate([]string{"@a", "", "(NEW REG)", "// only comment ", "D = M - 1 // comment"})
	if err == nil {
		fmt.Println(translatedRows)
	} else {
		fmt.Println(err)
	}
}
