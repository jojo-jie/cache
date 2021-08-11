package main

import (
	"fmt"
	"sync"
)

func main() {
	//runtime.GOMAXPROCS(8)
	/*trace.Start(os.Stderr)
	defer trace.Stop()*/
	urls := []string{"0.0.0.0:1121", "0.0.0.0:1122", "0.0.0.0:1123"}
	wg:=sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Println(url)
		}(url)
	}
	wg.Wait()
}
