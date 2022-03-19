package main

import (
	"fmt"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	// for a, b := range oqs.EnabledSigs() {
	// 	// fmt.Println(a)
	// 	fmt.Println("{2, 16, 840, 1, 101, 3, 4, 3, ", 40+a, "}", b)
	// }
	for _, b := range oqs.EnabledSigs() {
		fmt.Println(b)
	}
}
