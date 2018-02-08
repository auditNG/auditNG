package main

import (
	"fmt"
	"github.com/pre-processink/source"
	"github.com/pre-processink/transform"
)

func main() {
	var s source.Source = source.NewESSource()
	result, err := s.Fetch()

	if err != nil {
		fmt.Println("Error fetching from source")
		return
	}

	var t transform.Transform = transform.NewTransform()
	err = t.Process(result)

	if err != nil {
		fmt.Println("Error fetching from source")
		return
	}
}
