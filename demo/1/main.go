package main

/*
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"time"
)

func main() {
	C.srand(C.uint(time.Now().Unix()))

	randomNumer := C.rand()
	// var randomNumer C.int = C.rand()

	fmt.Printf("random number = %d\n", int(randomNumer))

	// pid := C.rand()
	var pid C.pid_t = C.getpid()

	fmt.Printf("my pid = %d\n", int(pid))
}
