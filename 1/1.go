package main

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
)

type silly struct {
	a 	*silly1
}

type silly1 struct {
}

func (s silly) MarshalEasyJSON(w *jwriter.Writer) {
	w.String("Hello")
}

func main() {
	var s *silly
	s = &silly{&silly1{}}
	data, err := easyjson.Marshal(s)
	fmt.Println(data, err)
}