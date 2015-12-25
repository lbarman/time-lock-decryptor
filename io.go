package main

import (
	"os"
	"encoding/json"
	"bufio"
	"fmt"
	"strings"
	"math/big"
	"io/ioutil"
)

func WriteBig(p *big.Int, file string) {
	data := p.Bytes()
	ioutil.WriteFile(file, data[:], 0644)
}

func ReadBig(file string) *big.Int {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("Can't read n.dat, please regenerate the primes.")
		panic(err.Error())
	}

	p := big.NewInt(0).SetBytes(data)
	return p
}


func JsonToString(data interface{}) string {
	b, err := json.Marshal(data)
    if err != nil {
        panic(err.Error())
        return ""
    }
    s := string(b)
    return s
}

func StringToJson(s string) *Puzzle {
	data := &Puzzle{}
	err := json.Unmarshal([]byte(s), data)
	if err != nil{
		panic(err.Error())
	}
	return data
}



func ReadBigFromConsole(text string) *big.Int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(text)
	s, _ := reader.ReadString('\n')
	s = strings.Replace(s, "\n", "", 1)

	y := big.NewInt(0)
	y.SetString(s, 10)

	return y
}