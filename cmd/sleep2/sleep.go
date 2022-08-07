package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	n := os.Args[1]
	a, _ := strconv.Atoi(n)
	time.Sleep(time.Duration(a) * time.Second)
}
