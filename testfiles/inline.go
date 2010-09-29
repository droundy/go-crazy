package main

import(
	"os"
	"fmt"
)

// The hello function just greets the world in a friendly way.
func hello(internal int) os.Error {
	for j:=0; j<internal; j++ {
		_,err := fmt.Println("Hello world!")
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	for i:=0; i<10; i++ {
		if e := hello(i); e != nil {
			fmt.Println(e)
			return
		}
	}
}
