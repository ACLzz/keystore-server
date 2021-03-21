package utils

import "os"

var SigCh = make(chan os.Signal, 1)
var EndCh = make(chan int, 1)
