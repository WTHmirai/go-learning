package main
import (
	"fmt"
	"strings"
)

func makeSuffix(suffix string) func(string) string {
    var n = 1
    defer fmt.Println("suffix = ", suffix, ", n = ", n)
    defer fmt.Println("...")
    n = n + 1
 
    fmt.Println("makeSuffix..")
    return func(file_name string) string {
        if (! strings.HasSuffix(file_name, suffix)) {
            return file_name + suffix
        }
                return file_name
    }
}
 
func main() {
    f := makeSuffix(".txt")
    fmt.Println(f("szc.txt"))
    fmt.Println(f("szc"))
}