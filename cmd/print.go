package main

import (
	"fmt"

	"github.com/a-h/templ"
)

func printDiff(b *templ.DiffBuffer) {
	b.Flush()
	fmt.Println("Segs:")
	for _, s := range b.Segs {
		fmt.Printf(" - '%v'\n", s)
	}
	fmt.Println("Values:")
	for _, v := range b.Vals {
		fmt.Printf(" - '%v'\n", v)
	}
}
