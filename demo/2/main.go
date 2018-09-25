package main

/*
#cgo CFLAGS:  -I${SRCDIR}/_mymath
#cgo LDFLAGS: -L${SRCDIR}/_mymath -lmymath

#include "mymath.h"
*/
import "C"

import (
	"fmt"
)

func main() {
	a := C.int(22)
	b := C.int(6)

	fmt.Printf("%d tambah %d adalah %d\n", int(a), int(b), int(C.tambah(a, b)))
	fmt.Printf("%d kurang %d adalah %d\n", int(a), int(b), int(C.kurang(a, b)))
	fmt.Printf("%d kali %d adalah %d\n", int(a), int(b), int(C.kali(a, b)))

	var rem C.int
	res := C.bagi(a, b, &rem)
	fmt.Printf("%d bagi %d adalah %d dengan sisa %d\n", int(a), int(b), int(res), int(rem))
}
