package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	inputStr := "LLRR=\n"
	expectStr := "input = output = 210122\n"

	stdinBackup := os.Stdin
	stdoutBackup := os.Stdout

	ri, wi, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	ro, wo, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = ri
	os.Stdout = wo

	go func(w *os.File) {
		time.Sleep(time.Microsecond * 100)
		fmt.Fprint(w, inputStr)
		w.Close()
	}(wi)
	main()

	wo.Close()
	o, err := io.ReadAll(ro)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdin = stdinBackup
	os.Stdout = stdoutBackup

	got := string(o)
	want := expectStr
	if got != want {
		t.Errorf("main() = %q, want %q", got, want)
	}
}

func TestLreq(t *testing.T) {

	type TestCase struct {
		testName    string
		inputEnc    string
		expectDec   string
		expectPanic bool
	}

	testCases := []TestCase{
		{testName: "TestLreq LLRR=", inputEnc: "LLRR=", expectDec: "210122", expectPanic: false},
		{testName: "TestLreq LRLL", inputEnc: "LRLL", expectDec: "21210", expectPanic: false},
		{testName: "TestLreq LRRR=", inputEnc: "LRRR=", expectDec: "101233", expectPanic: false},
		{testName: "TestLreq LLRLR=", inputEnc: "LLRLR=", expectDec: "2101011", expectPanic: false},
		{testName: "TestLreq RRL", inputEnc: "=RRL=", expectDec: "001200", expectPanic: false},
		{testName: "TestLreq RRLL", inputEnc: "RRLL", expectDec: "01210", expectPanic: false},
		{testName: "TestLreq RRLLL", inputEnc: "RRLLL", expectDec: "023210", expectPanic: false},
		{testName: "TestLreq RRLLL", inputEnc: "RRLLLL", expectDec: "0343210", expectPanic: false},
		{testName: "TestLreq RR=LL=", inputEnc: "RR=LL=", expectDec: "0122100", expectPanic: false},
		{testName: "TestLreq RR=LL=L", inputEnc: "RR=LL=L", expectDec: "02332110", expectPanic: false},
		{testName: "TestLreq LRRRLLL", inputEnc: "LRRRLLL", expectDec: "10123210", expectPanic: false},
		{testName: "TestLreq LRRRLLLL", inputEnc: "LRRRLLLL", expectDec: "212343210", expectPanic: false},
		{testName: "TestLreq LRRRLLL", inputEnc: "RRRLLL", expectDec: "0123210", expectPanic: false},
		{testName: "TestLreq panic encoding", inputEnc: "AEIOU", expectDec: "", expectPanic: true},
		{testName: "TestLreq panic encoding", inputEnc: "RRRRRRRRRR", expectDec: "", expectPanic: true},
		{testName: "TestLreq panic encoding", inputEnc: "LLLLLLLLLL", expectDec: "", expectPanic: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			if testCase.expectPanic == true {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic")
					}
				}()
			}
			actualDec := lreq(testCase.inputEnc)

			assert.Equal(t, testCase.expectDec, actualDec)
		})
	}
}

func TestLprocess(t *testing.T) {
	type TestCases struct {
		testName  string
		input     string
		expectDec string
		err       error
	}

	testCases := []TestCases{
		{testName: "TestLprocess 1", input: "1", expectDec: "10", err: nil},
		{testName: "TestLprocess 0", input: "0", expectDec: "10", err: nil},
		{testName: "TestLprocess 10", input: "10", expectDec: "210", err: nil},
		{testName: "TestLprocess 01", input: "01", expectDec: "010", err: nil},
		{testName: "TestLprocess 010", input: "010", expectDec: "1210", err: nil},
		{testName: "TestLprocess 1210", input: "1210", expectDec: "23210", err: nil},
		{testName: "TestLprocess 19210", input: "19210", expectDec: "", err: errors.New("mulform")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			actualDec, err := lprocess(testCase.input)

			if testCase.err != nil {
				assert.Error(t, testCase.err, err.Error())
				return
			}
			assert.Equal(t, testCase.expectDec, actualDec)
		})
	}

}

func TestRprocess(t *testing.T) {
	type TestCases struct {
		testName  string
		input     string
		expectDec string
		err       error
	}

	testCases := []TestCases{
		{testName: "TestLprocess 1", input: "1", expectDec: "12", err: nil},
		{testName: "TestLprocess 0", input: "0", expectDec: "01", err: nil},
		{testName: "TestLprocess 10", input: "10", expectDec: "101", err: nil},
		{testName: "TestLprocess 12345678", input: "12345678", expectDec: "123456789", err: nil},
		{testName: "TestLprocess 123456789", input: "123456789", expectDec: "", err: errors.New("mulform")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			actualDec, err := rprocess(testCase.input)

			if testCase.err != nil {
				assert.Error(t, testCase.err, err.Error())
				return
			}
			assert.Equal(t, testCase.expectDec, actualDec)
		})
	}
}

func TestLplus(t *testing.T) {
	type TestCases struct {
		testName  string
		input     string
		expectDec string
		err       error
	}

	testCases := []TestCases{
		{testName: "TestLplus 876543210", input: "876543210", expectDec: "9876543210", err: nil},
		{testName: "TestLplus 12345321", input: "12345321", expectDec: "123454320", err: nil},
		{testName: "TestLplus 1234598", input: "1234598", expectDec: "", err: errors.New("lplus mulform")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			out, err := lplus(testCase.input)

			if testCase.err != nil {
				assert.Error(t, testCase.err, err.Error())
				return
			}
			assert.Equal(t, testCase.expectDec, out)
		})
	}

}
