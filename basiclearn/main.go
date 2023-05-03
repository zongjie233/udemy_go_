package main

import "fmt"

func main() {
	fmt.Println("hello world")

	var whatToSay string

	var i int

	whatToSay = "goodbye, cool world"
	fmt.Println(whatToSay)

	i = 7
	fmt.Println("i is set to", i)
	whatWasSaid, theOtherSaid := saySomething()
	fmt.Println("func returned", whatWasSaid, theOtherSaid)

}

func saySomething() (string, string) {
	return "something", "else"
}
