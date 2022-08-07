package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("sleep")
	go func() {
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(1 * time.Second)
	cmd.Process.Kill()

	time.Sleep(1 * time.Second)
	fmt.Println(cmd.ProcessState.ExitCode())
}
