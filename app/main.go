package main

import "fmt"

func main() {
	b := make(chan []byte)
	isC := make(chan struct{})
	bytes := []byte{'1', '1', '1', '1', '1', '1', '1'}

	go func() {
		for i := 0; i < len(bytes); i++ {
			bt := make([]byte, 1)
			bt[0] = bytes[i]
			b <- bt
		}
		close(b)
	}()

	go func() {
		for x := range b {
			fmt.Println(x)
		}
		isC <- struct{}{}
	}()

	<-isC
}
