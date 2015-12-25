package main

import (
	"os"
	"flag"
	"bufio"
	"fmt"
	"math/rand"
	"time"
	"strings"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

var expOneSecond *big.Int

func newPrime(pow int) *big.Int {
	//we pick the first prime
	base      := int64(2)
	power     := int64(pow)
	init      := ExpByPowOfTwo(big.NewInt(base), big.NewInt(power))
	precision := 50
	one       := big.NewInt(1)
	inc       := 0
	for !init.ProbablyPrime(precision) {
		init = big.NewInt(0).Add(init, one)
		inc += 1
	}
	prime := init
	fmt.Println("Parameters are :")
	fmt.Println("Prime number is", prime)
	fmt.Println("Corresponding to", base, "^", power, "+", inc)
	fmt.Println("Its hash is", sha256AndHex(prime.Bytes()))

	return prime
}

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

func main() {
	fmt.Println("Time-Lock Puzzle")
	fmt.Println("****************")

	expOneSecond =ExpByPowOfTwo(big.NewInt(2), big.NewInt(5000))

	gen := flag.Bool("gen", false, "Start generating a new puzzle")
	solve := flag.Bool("solve", true, "Start solving a puzzle")

	if *gen {
		preparePuzzle()
	} else if *solve {
		solvePuzzle()
	} else {
		fmt.Println("Please add the flag -gen or -solve in the command line for the desired action.")
	}

	fmt.Println("****************")
}

func preparePuzzle() {

	//randomness
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd       := rand.New(rndSource)

	//we generate the parameters
	p := newPrime(1234)
	q := newPrime(1010)
	one := big.NewInt(1)
	p_1 := big.NewInt(0).Sub(p, one)
	q_1 := big.NewInt(0).Sub(q, one)
	phi := big.NewInt(0).Mul(p_1, q_1)
	n := big.NewInt(0).Mul(p, q)

	fmt.Println("Parameters are :")
	fmt.Println("N is", n)
	fmt.Println("Its hash is", sha256AndHex(n.Bytes()))

	//generate benchmark, a 1-second exponentiation
	repeat := 10
	times := make([]time.Duration, repeat)

	for i:=0; i<repeat; i++ {
		rand      := big.NewInt(0).Rand(rnd, n)
		start     := time.Now()

		ExpByPowOfTwoModular(rand, expOneSecond, n)

		diff := time.Now().Sub(start)
		times[i] = diff
		fmt.Println("Benchmark", i, "took", diff)
	}

	fmt.Println("****************")
	fmt.Println("All clear !")

	proceed := false
	x := big.NewInt(0)
	y := big.NewInt(0)

	for !proceed{
		x = big.NewInt(0).Rand(rnd, n)

		y = readBig("Enter number of iterations : ")

	    fmt.Println("Iter is", y)
	    duration := time.Duration(int(y.Int64())) * time.Second
	    fmt.Println("Expected duration is", duration)
	    fmt.Println("Proceed [y/n] ?")
	    reader := bufio.NewReader(os.Stdin)
	    ans, _ := reader.ReadString('\n')
	    ans = strings.Replace(ans, "\n", "", 1)

	    if ans == "y" {
	    	proceed = true
	    }
	}

	curr    := big.NewInt(0)
	zero    := big.NewInt(0)
	hundred := big.NewInt(100)
	base    := x
	iter    := y

	exponent := one
	for curr.Cmp(y) < 0 {
		if mod(curr, hundred).Cmp(zero) == 0 {
			diff := big.NewInt(0).Sub(y, curr)
			rem := time.Duration(int(diff.Int64())) * time.Second
			fmt.Println("Iteration", curr.String()+"/"+y.String(), "remaning time is", rem)
		}
		
		exponent = mod(big.NewInt(0).Mul(exponent, expOneSecond), phi)

		curr.Add(curr,one)

	}

	x = ExpByPowOfTwoModular(x, exponent, n)

	fmt.Println("################")
	fmt.Println("Time-Lock Puzzle finished !")
	fmt.Println("")
	fmt.Println("time to unlock:")
	fmt.Println(time.Duration(int(y.Int64()))*time.Second)
	fmt.Println("N:")
	fmt.Println(n)
	fmt.Println("base:")
	fmt.Println(base)
	fmt.Println("expOneSecond:")
	fmt.Println(expOneSecond)
	fmt.Println("iter:")
	fmt.Println(iter)
	fmt.Println("result:")
	fmt.Println(x)
	fmt.Println("sha256 of result:")
	fmt.Println(sha256AndHex(x.Bytes()))
	fmt.Println("")
	fmt.Println("################")
}

func solvePuzzle() {
	//we generate the parameters
	n := readBig("Please enter n, the modulus :")
	fmt.Println("n's hash is", sha256AndHex(n.Bytes()))
	base := readBig("Please enter the base :")
	fmt.Println("base's hash is", sha256AndHex(base.Bytes()))
	iter := readBig("Please enter the number of iterations :")

	x := base

	curr := big.NewInt(0)
	one := big.NewInt(1)
	zero := big.NewInt(0)
	hundred := big.NewInt(100)

	for curr.Cmp(iter) < 0 {
		if mod(curr, hundred).Cmp(zero) == 0 {
			diff := big.NewInt(0).Sub(iter, curr)
			rem := time.Duration(int(diff.Int64())) * time.Second
			fmt.Println("Iteration", curr.String()+"/"+iter.String(), ", remaining time is", rem)
		}
		
		x = ExpByPowOfTwoModular(x, expOneSecond, n)

		curr.Add(curr,one)
	}

	fmt.Println("################")
	fmt.Println("Time-Lock Puzzle solved !")
	fmt.Println("")
	fmt.Println(x)
	fmt.Println("")
	fmt.Println(sha256AndHex(x.Bytes()))
	fmt.Println("")
	fmt.Println("################")
}

func readBig(text string) *big.Int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(text)
	s, _ := reader.ReadString('\n')
	s = strings.Replace(s, "\n", "", 1)

	y := big.NewInt(0)
	y.SetString(s, 10)

	return y
}

func sha256AndHex(x []byte) string {
	hash := sha256.Sum256(x)
	return hex.EncodeToString(hash[:])
}

func ExpByPowOfTwoModular(base, power, modulus *big.Int) *big.Int {

	//if power == 0, return 1
	if power.Cmp(big.NewInt(0)) == 0 {
		return mod(big.NewInt(1), modulus)
	}
	//if power == 1,  base
	if power.Cmp(big.NewInt(1)) == 0 {
		return mod(base, modulus)
	}

    result := big.NewInt(1)
    one := big.NewInt(1)
    for power.Cmp(one) > 0 {
        if modBy2(power).Cmp(one) == 0 {
            result = mod(multiply(result, base), modulus)
        	power = big.NewInt(0).Sub(power, one)
        }
        power = divideBy2(power)
        base = mod(multiply(base, base), modulus)
    }
    return mod(big.NewInt(0).Mul(result, base), modulus)
}


func ExpByPowOfTwo(base, power *big.Int) *big.Int {

	//if power == 0, return 1
	if power.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1)
	}
	//if power == 1,  base
	if power.Cmp(big.NewInt(1)) == 0 {
		return base
	}

    result := big.NewInt(1)
    one := big.NewInt(1)
    for power.Cmp(one) > 0 {
        if modBy2(power).Cmp(one) == 0 {
            result = multiply(result, base)
        	power = big.NewInt(0).Sub(power, one)
        }
        power = divideBy2(power)
        base = multiply(base, base)
    }
    return big.NewInt(0).Mul(result, base)
}

func modBy2(x *big.Int) *big.Int {
    return big.NewInt(0).Mod(x, big.NewInt(2))
}

func divideBy2(x *big.Int) *big.Int {
    return big.NewInt(0).Div(x, big.NewInt(2))
}

func multiply(x, y *big.Int) *big.Int {
    return big.NewInt(0).Mul(x, y)
}

func mod(x, m *big.Int) *big.Int {
	return big.NewInt(0).Mod(x, m)
}