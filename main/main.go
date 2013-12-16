package main

import (
	"fmt"
	"github.com/creack/btce"
	"log"
)

func main() {
	b := &btce.Api{
		Url: "https://btc-e.com/tapi",
	}
	// ok, err := b.GetInfo()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//	log.Printf("%#v\n", ok)

	options := &btce.Options{
		Count: 1,
	}

	ok2, err := b.TransHistory(options)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range ok2 {
		fmt.Printf("\n\n%s\n%#v\n", k, v)
	}
	log.Printf("%#v\n", ok2)

	fmt.Printf("len: %d\n", len(ok2))
	println("\nOK")
}
