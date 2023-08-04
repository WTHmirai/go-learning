package main
import "fmt"

func makeSuffix(suffix string) (func(string) string) {
	return func(filename string) string {
		if strings.HasSuffix(filename,suffix) {
			return filename
		} else {
			return filename + suffix
		}
	}
}

func main() {
	f := makeSuffix(".jpg")
	fmt.Println(f("wth.jpg"))
	fmt.Println(f("wthasfasf"))
}