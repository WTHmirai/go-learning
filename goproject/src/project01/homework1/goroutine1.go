package main
import (
	"fmt"
)

func Writedata(numchan chan int) {
	for i:=1 ; i<=2000; i++ {
		numchan <- i
	}

	close(numchan)
}

func Readdata(numchan chan int,reschan chan int,exitchan chan bool) {
	for {
		v,ok := <- numchan
		if !ok {
			break
		}
		reschan <- ((1+v)*v/2)
	}
	exitchan <- true
}

func main() {
	numchan := make(chan int,2000)
	reschan := make(chan int,2000)
	exitchan := make(chan bool,8)

	go Writedata(numchan)
	for i:=0; i<8; i++ {
		go Readdata(numchan,reschan,exitchan)
	}

	flag := 0
	for {
		if flag == 8 {
			fmt.Println("从这里出去的")
			break
		}

		v,ok := <- exitchan
		if !ok {
			break
		}

		if v {
			flag += 1
		}
	}
	close(reschan)

	for {
		v,ok := <- reschan
		if !ok {
			break
		}
		fmt.Println(v)
	}

	close(exitchan)

}