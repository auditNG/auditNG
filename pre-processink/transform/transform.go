package transform

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"strconv"
	"strings"
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
	sysstr, err := t.getStringValue(message, "syscall=")
	if nil != err {
		fmt.Println("Unable to get syscall")
	}
	fmt.Println("syscall: " + sysstr)

	syscall, err := t.getIntValue(message, "syscall=")
	if nil != err {
		fmt.Println("Unable to get syscall")
	}
	fmt.Println("syscall: " + strconv.Itoa(syscall))

	exitcode, err := t.getIntValue(message, "exit=")
	if nil != err {
		fmt.Println("Unable to get exitcode")
	}
	fmt.Println("exitcode: " + strconv.Itoa(exitcode))

}

func (t Transform) getStringValue(message string, key string) (string, error) {
	data := message
	start := 0
	end := 0

	if start = strings.Index(data, key); start < 0 {
		return "", errors.New("Error parsing exit code")
	}

	// Progress the start point beyond the = sign
	start += len(key)
	if end = strings.IndexByte(data[start:], spaceChar); end < 0 {
		// There was no ending space, maybe the syscall id is at the end of the line
		end = len(data) - start
	}

	retval := data[start : start+end]
	return retval, nil
}

func (t Transform) getIntValue(message string, key string) (int, error) {
	val, err := t.getStringValue(message, key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}
