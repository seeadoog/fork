## Fork
提供一个快速fork 子进程并管理子进程的工具,支持父进程和子进程相互通信


## Example

````go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/seeadoog/fork"
	"os"
	"runtime"
)

type Data struct {
	Num int
}

func main() {
	f := fork.NewForker(2)
	outputs := make([]*bytes.Buffer, 0)
	f.SetPreForkChild(func(child *fork.Cmd) error {
		bf := bytes.NewBuffer(nil)
		outputs = append(outputs, bf)
		child.Cmd().Stdout = bf
		return nil
	})

	err := f.ForkProcess(func(f *fork.MasterTool) error {
		return nil
	}, func(c *fork.ChildrenTool) error {
		runtime.GOMAXPROCS(1)
		return json.NewEncoder(os.Stdout).Encode(&Data{
			Num: 5,
		})
	})
	if !f.IsMaster() {
		return
	}
	if err != nil {
		panic(err)
	}
	f.Wait()
	sum := 0
	for _, output := range outputs {
		data := new(Data)
		err := json.NewDecoder(output).Decode(data)
		if err != nil {
			panic(err)
		}
		sum += data.Num
	}
	fmt.Println("sum is :", sum)
}

````