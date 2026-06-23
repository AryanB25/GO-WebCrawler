package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://gobyexample.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Status: ", resp.Status)

	bytes, err := io.ReadAll(resp.Body) 
	text := string(bytes)
	original_slice := text[:200]
	fmt.Println(original_slice)
}
