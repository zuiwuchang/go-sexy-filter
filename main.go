package main

import (
	"log"
	"net/url"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	proxy, e := url.Parse("http://127.0.0.1:8118")
	if e != nil {
		log.Fatalln(e)
	}

	t2, e := NewT66y2()
	if e != nil {
		log.Fatal(e)
	}
	t2.SetProxy(proxy)
	for i := 0; i < 10; i++ {
		t2.Get(i)
	}
}
