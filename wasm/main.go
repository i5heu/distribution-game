package main

//export hostPrint
func hostPrint(msg string)

func main() {
	hostPrint("Hello, world!")
}

//export add_one
func add_one(x int) (int, string) {
	hostPrint("Hello, world!")
	return x * x, "Hello"
}
