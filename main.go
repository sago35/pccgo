package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/sago35/pccgo/exec"
)

var (
	ch    = make(chan *exec.Cmd, 10)
	limit = make(chan struct{}, 4)
	wg    sync.WaitGroup
)

func main() {
	runServer()
}

func runServer() {
	go buildLoop()

	fmt.Println("started.")
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:9876", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rb, _ := ioutil.ReadAll(r.Body)

	var e exec.Cmd
	err := json.Unmarshal(rb, &e)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	fmt.Println(e.Target)
	if e.Join {
		// Join時は、並列実行されたJobの終了を待つ
		wg.Wait()
		e.Run()
	} else {
		ch <- &e
	}
	fmt.Fprintf(w, "done")
}

func buildLoop() {
	var empty struct{}

	for {
		select {
		case limit <- empty:
			e := <-ch
			wg.Add(1)
			go func(x *exec.Cmd) {
				x.Run()
				wg.Done()
				<-limit
			}(e)
		}
	}
}
