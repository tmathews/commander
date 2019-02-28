# Commander

A small utility for executing N sub commands with flags.

## Example

```go
package main

import (
	"fmt"
	"os"
	"flag"
	
	"github.com/tmathews/commander"
)

func main() {
	var args []string
	if len(os.Args) < 2 {
		args = []string{}
	} else {
		args = os.Args[1:]
	}
	err := commander.Exec(commander.M{
		"foo":      foo,
		"version":  ver,
	}, args)
	if err != nil {
		switch v := err.(type) {
		case commander.Error:
			fmt.Print(v.Help())
			os.Exit(2)
		default:
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func ver(name string, args []string) error {
	fmt.Println("1.0")
	return nil
}

func foo(name string, args []string) error {
	return commander.Exec(commander.M{
		"bar": func(name string, args []string) error {
			var i int
			set := flag.NewFlagSet(name, flag.ExitOnError)
			set.IntVar(&i, "i", 0, "an integer")
			if err := set.Parse(args); err != nil {
				return err
			}
			fmt.Printf("foo bar %d\n", i)
			return nil
		},
	}, args)
}
```