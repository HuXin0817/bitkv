package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/HuXin0817/bitkv/api"
	"github.com/zeromicro/go-zero/core/logx"
)

var serverAddr = flag.String("h", "127.0.0.1:7070", "server address")

var client *api.Client

func argsErr(length, legalLength int) error {
	if length == legalLength {
		return nil
	}
	return fmt.Errorf("error args number, except %d, but %d", legalLength, length)
}

var funcMap = map[string]func(args []string) error{

	// PUT <key> <value>
	"PUT": func(args []string) error {
		if err := argsErr(len(args), 3); err != nil {
			return err
		}

		if err := client.Put(args[1], args[2]); err != nil {
			return err
		}

		fmt.Println("done.")
		return nil
	},

	// GET <key>
	"GET": func(args []string) error {
		if err := argsErr(len(args), 2); err != nil {
			return err
		}

		r, err := client.Get(args[1])
		if err != nil {
			return err
		}

		if r != "" {
			fmt.Println(r)
		}

		return nil
	},

	// DELETE <key>
	"DELETE": func(args []string) error {
		if err := argsErr(len(args), 2); err != nil {
			return err
		}

		if err := client.Delete(args[1]); err != nil {
			return err
		}

		fmt.Println("done.")
		return nil
	},

	// EXIT
	"EXIT": func(args []string) error {
		if err := argsErr(len(args), 1); err != nil {
			return err
		}

		fmt.Println("bye.")
		os.Exit(0)
		return nil
	},
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(137)
		}
	}()

	flag.Parse()
	client = api.NewClient(*serverAddr)
	scanner := bufio.NewScanner(os.Stdin)
	logx.Disable()

	for {
		fmt.Printf("bitkv@%s> ", *serverAddr)
		scanner.Scan()
		args := strings.Fields(scanner.Text())
		if len(args) == 0 {
			continue
		}

		if f, ok := funcMap[strings.ToUpper(args[0])]; ok {
			if err := f(args); err != nil {
				fmt.Printf("BITKV error: %s\n", err)
			}
		}
	}
}
