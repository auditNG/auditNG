package transform

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
  "strings"
  "strconv"
)

var spaceChar = byte(' ')

func NewTransform() Transform {
	return Transform{}
}

type Transform struct {
}

func (t Transform) Process(input string) error {

	jsonparser.ArrayEach([]byte(input),
		func(actVal []byte, _ jsonparser.ValueType, _ int, err error) {
			timestamp, err := jsonparser.GetString(actVal, "_source", "@timestamp")
			if err != nil {
				fmt.Println("JSON parsing error: ", err)
				return
			}
			fmt.Println("timestamp: " + timestamp)

			message, err := jsonparser.GetString(actVal, "_source", "message")
			if err != nil {
				fmt.Println("JSON parsing error: ", err)
				return
			}
			fmt.Println("message: " + message)
      t.processMessage(message)

		}, "hits", "hits")

	return nil
}

func (t Transform) processMessage(message string) {
	syscall, err := t.getSysCall(message)
	if nil != err {
		fmt.Println("Unable to get syscall")
	}

	fmt.Println("syscall is:" + string(syscall))
}

func (t Transform) getSysCall(message string) (int, error) {
	data := message
	start := 0
	end := 0

	if start = strings.Index(data, "syscall="); start < 0 {
		return 0, errors.New("Error parsing syscall")
	}

	// Progress the start point beyond the = sign
	start += 8
	if end = strings.IndexByte(data[start:], spaceChar); end < 0 {
		// There was no ending space, maybe the syscall id is at the end of the line
		end = len(data) - start

		// If the end of the line is greater than 5 characters away (overflows a 16 bit uint) then it can't be a syscall id
		if end > 5 {
      return 0, errors.New("Error parsing syscall")
		}
	}

	syscall := data[start : start+end]
  fmt.Println("syscall is:" + syscall)
	return strconv.Atoi(syscall)
}
