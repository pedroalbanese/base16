package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"encoding/hex"
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
			// Hex dump
			encoded := hex.Dump([]byte(inputData))
			fmt.Println(encoded)
		} else if *dec {
			// Decodificar a partir do hexadecimal
			decoded, err := hex.DecodeString(inputData)
			if err != nil {
				fmt.Println("Error decoding from hexadecimal:", err)
				os.Exit(1)
			}
			fmt.Println(string(decoded))
		} else {
			// Codificar para hexadecimal
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
				// Hex dump
				encoded := hex.Dump([]byte(inputData))
				fmt.Println(encoded)
			} else if *dec {
				// Decodificar a partir do hexadecimal
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
				// Codificar para hexadecimal e dividir em colunas
				encoded := hex.EncodeToString([]byte(inputData))
				for _, chunk := range split(encoded, *col) {
					fmt.Println(chunk)
				}
			}
		} else {
			if *dec {
				// Decodificar a partir do hexadecimal
				decoded, err := hex.DecodeString(inputData)
				if err != nil {
					fmt.Println("Error decoding from hexadecimal:", err)
					os.Exit(1)
				}
				fmt.Println(string(decoded))
			} else {
				// Codificar para hexadecimal
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

	// Split the lines of the hexdump
	lines := strings.Split(input, "\n")

	// Iterate through the lines and collect hexadecimal characters
	for _, line := range lines {
		// Ignore lines with less than 59 characters
		if len(line) < 59 {
			continue
		}

		// Extract characters from column 10 to 58
		hexCharsInLine := line[9:58]

		// Remove spaces
		hexCharsInLine = strings.ReplaceAll(hexCharsInLine, " ", "")

		// Append cleaned hex characters to the buffer
		buffer.WriteString(hexCharsInLine)
	}

	// Decode the filtered hexadecimal characters
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