package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sago35/pccgo/exec"

	"fmt"
	"os"
	"strings"
)

func main() {
	var ec exec.Cmd

	if strings.HasSuffix(os.Args[0], ".join.exe") {
		ec.Path = strings.Replace(os.Args[0], ".join.exe", ".exe", 1)
		ec.Join = true
	} else if strings.HasSuffix(os.Args[0], ".wo_join.exe") {
		ec.Path = strings.Replace(os.Args[0], ".wo_join.exe", ".exe", 1)
		ec.Join = false
	} else {
		fmt.Fprintf(os.Stderr, "exe not found")
		os.Exit(1)
	}

	args := os.Args[1:]
	for i, arg := range args {
		ec.Args = append(ec.Args, arg)
		if arg == "-o" {
			fmt.Println(args[i+1])
			ec.Target = args[i+1]
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	ec.Dir = wd
	ec.Env = os.Environ()

	if false {
		ec.Stdout = os.Stdout
		ec.Stderr = os.Stderr
		exitStatus := ec.Run()
		if exitStatus != 0 {
			fmt.Fprintf(os.Stderr, "error %d\n", exitStatus)
		}
		os.Exit(exitStatus)
	} else {
		j, err := json.Marshal(ec)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		resp, err := http.Post("http://127.0.0.1:9876", "application/json", bytes.NewBuffer(j))
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintln(os.Stderr, body)
			os.Exit(1)
		}

		if ec.Join {
			os.Exit(0)
		} else {
			os.Exit(0)
		}
	}
}
