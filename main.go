package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Fprintf(os.Stdout, "input = ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.ToUpper(scanner.Text())

	fmt.Fprintf(os.Stdout, "output = %s\n", lreq(input))
}

func lreq(enc string) string {
	var err error
	output := "0"

	for i, v := range enc {
		switch v {
		case 'L':
			output, err = lprocess(output)
			if err != nil {
				panic(err)
			}
		case 'R':
			output, err = rprocess(output)
			if err != nil {
				panic(err)
			}
		case '=':
			output += string(output[len(output)-1])
		default:
			panic("invalid encoding")
		}

		if i == len(enc)-1 {
			if enc[0] == 'R' && output[0] != '0' {
				runes := []rune(output)
				runes[0] = '0'
				output = string(runes)
			}
		}
	}

	return output
}

func lprocess(in string) (dec string, err error) {
	if in[len(in)-1] != '0' {
		return in + "0", nil
	}

	if strings.Contains(in, "9") {
		return "", errors.New("mulform")
	}

	return lplus(in)
}

func rprocess(in string) (dec string, err error) {
	if in[len(in)-1] == '9' {
		return "", errors.New("mulform")
	}
	return in + string(in[len(in)-1]+1), nil
}

func lplus(str string) (string, error) {
	runes := []rune(str)

	for i := len(runes) - 1; i >= 0; i-- {
		digit := int(runes[i] - '0')

		if digit == 9 {
			return "", errors.New("lplus mulform")
		} else {
			digit = digit + 1
		}

		runes[i] = rune(digit + '0')

		if i == 0 {
			break
		} else if runes[i] < runes[i-1] {
			break
		}
	}
	return string(runes) + "0", nil
}
