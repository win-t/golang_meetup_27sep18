package main

/*
#include "lala.h"
*/
import "C"

import (
	"fmt"
)

func main() {
	C.lala()
}

//export go_tambah
func go_tambah(a, b C.int) C.int {
	fmt.Printf("inside go_tambah(%d, %d)\n", int(a), int(b))
	return C.int(int(a) + int(b))
}
