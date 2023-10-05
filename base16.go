package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	col  = flag.Int("w", 64, "Wrap lines after N columns")
	dec  = flag.Bool("d", false, "Decode instead of Encode")
	dump = flag.Bool("dump", false, "Hexdump instead of Encode/Decode")
)

func main() {
	flag.Parse()

	if *col == 0 && len(flag.Args()) > 0 {
		inputFile := flag.Arg(0)

		data, err := ioutil.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading the file:", err)
			os.Exit(1)
		}

		inputData := string(data)
		inputData = strings.TrimSuffix(inputData, "\r\n")
		inputData = strings.TrimSuffix(inputData, "\n")

		if *dump {
			encoded := hex.Dump([]byte(inputData))
			fmt.Println(encoded)
		} else if *dec {
			decoded, err := hex.DecodeString(inputData)
			if err != nil {
				fmt.Println("Error decoding from hexadecimal:", err)
				os.Exit(1)
			}
			fmt.Println(string(decoded))
		} else {
			encoded := hex.EncodeToString([]byte(inputData))
			fmt.Println(encoded)
		}
	} else {
		var inputData string

		if len(flag.Args()) == 0 {
			data, _ := ioutil.ReadAll(os.Stdin)
			inputData = string(data)
		} else {
			inputFile := flag.Arg(0)

			data, err := ioutil.ReadFile(inputFile)
			if err != nil {
				fmt.Println("Error reading the file:", err)
				os.Exit(1)
			}
			inputData = string(data)
		}

		inputData = strings.TrimSuffix(inputData, "\r\n")
		inputData = strings.TrimSuffix(inputData, "\n")

		if *col != 0 {
			if *dump {
				encoded := hex.Dump([]byte(inputData))
				fmt.Println(encoded)
			} else if *dec {
				var decoded []byte
				var err error
				if !isHexDump(inputData) {
					decoded, err = decodeHexDump(inputData)
				} else {
					inputData = strings.ReplaceAll(inputData, "\r\n", "")
					inputData = strings.ReplaceAll(inputData, "\n", "")
					decoded, err = hex.DecodeString(inputData)
				}
				if err != nil {
					fmt.Println("Error decoding from hexadecimal:", err)
					os.Exit(1)
				}
				fmt.Println(string(decoded))
			} else {
				encoded := hex.EncodeToString([]byte(inputData))
				for _, chunk := range split(encoded, *col) {
					fmt.Println(chunk)
				}
			}
		} else {
			if *dec {
				decoded, err := hex.DecodeString(inputData)
				if err != nil {
					fmt.Println("Error decoding from hexadecimal:", err)
					os.Exit(1)
				}
				fmt.Println(string(decoded))
			} else {
				encoded := hex.EncodeToString([]byte(inputData))
				fmt.Println(encoded)
			}
		}
	}
}

func split(s string, size int) []string {
	ss := make([]string, 0, len(s)/size+1)
	for len(s) > 0 {
		if len(s) < size {
			size = len(s)
		}
		ss, s = append(ss, s[:size]), s[size:]
	}
	return ss
}

func decodeHexDump(input string) ([]byte, error) {
	var decoded []byte
	var buffer bytes.Buffer

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if len(line) < 59 {
			continue
		}

		hexCharsInLine := line[9:58]

		hexCharsInLine = strings.ReplaceAll(hexCharsInLine, " ", "")

		buffer.WriteString(hexCharsInLine)
	}

	decoded, err := hex.DecodeString(buffer.String())
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func isHexDump(input string) bool {
	if strings.Contains(input, "|") {
		return false
	} else {
		return true
	}
}
