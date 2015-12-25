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
	/*
	fmt.Println("Parameters are :")
	fmt.Println("Prime number is", prime)
	fmt.Println("Corresponding to", base, "^", power, "+", inc)
	fmt.Println("Its hash is", sha256AndHex(prime.Bytes()))
	*/

	return prime
}

func main() {
	fmt.Println("Time-Lock Puzzle")
	fmt.Println("****************")

	expOneSecond =ExpByPowOfTwo(big.NewInt(2), big.NewInt(50000))

	gen := flag.Bool("gen", true, "Start generating a new puzzle")
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

	/*
	fmt.Println("Parameters are :")
	fmt.Println("N is", n)
	fmt.Println("Its hash is", sha256AndHex(n.Bytes()))
	*/
	fmt.Println("Generation of parameters done...")
	fmt.Println("Benchmarking, each step should take around 1 sec.")

	//generate benchmark, a 1-second exponentiation
	repeat := 10
	sum := int64(0)

	for i:=0; i<repeat; i++ {
		rand      := big.NewInt(0).Rand(rnd, n)
		start     := time.Now()

		ExpByPowOfTwoModular(rand, expOneSecond, n)

		diff := time.Now().Sub(start)
		sum += diff.Nanoseconds()
		fmt.Println("Benchmark", i, "...")
	}

	mean := time.Duration(sum / 10) * time.Nanosecond
	fmt.Println("On this machine, one cycle takes ", mean)
	fmt.Println("Gonna generate the puzzle...")

	proceed := false
	x := big.NewInt(0)
	y := big.NewInt(0)

	for !proceed{
		x = big.NewInt(0).Rand(rnd, n)

		y = ReadBigFromConsole("Enter number of cycles (keep in mind that on other machines, one cycle might go faster or slower) : ")

	    duration := time.Duration(int(y.Int64())) * time.Second
	    fmt.Println("You entered ", y, ", expected duration is", duration, ". Proceed [y/n] ?")
 	    reader := bufio.NewReader(os.Stdin)
	    ans, _ := reader.ReadString('\n')
	    ans = strings.Replace(ans, "\n", "", 1)

	    if ans == "y" || ans == "Y" {
	    	proceed = true
	    }
	}

	curr    := big.NewInt(0)
	base    := x
	iter    := y

	fmt.Println("Generating...")

	exponent := one
	for curr.Cmp(y) < 0 {
		exponent = mod(big.NewInt(0).Mul(exponent, expOneSecond), phi)
		curr.Add(curr,one)

	}

	x = ExpByPowOfTwoModular(x, exponent, n)
	timeToUnlock := time.Duration(int(y.Int64()))*time.Second
	sha := sha256AndHex(x.Bytes())
	shaOfSha := sha256AndHex([]byte(sha))

	fmt.Println("Time-Lock Puzzle finished ! the solution is : ")
	fmt.Println("")
	fmt.Println(sha)
	fmt.Println("")
	fmt.Println("Please save the following JSON to recompute the solution (expected time :", timeToUnlock, ")")

	filledPuzzle := &Puzzle{n.String(), timeToUnlock.String(), iter.String(), base.String(), expOneSecond.String(), shaOfSha}

	fmt.Println("")
	fmt.Println(JsonToString(filledPuzzle))
	fmt.Println("")
}


func solvePuzzle() {

	fmt.Println("Please paste the JSON puzzle :")

	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans    = strings.Replace(ans, "\n", "", 1)

	puzzle := StringToJson(ans)

	//we generate the parameters
	n   , _ := big.NewInt(0).SetString(puzzle.N, 10)
	base, _ := big.NewInt(0).SetString(puzzle.Base, 10)
	iter, _ := big.NewInt(0).SetString(puzzle.NCycles, 10)
	exp , _ := big.NewInt(0).SetString(puzzle.ExponentOneSecond, 10)

	fmt.Println("All right ! Going to solve...")
	fmt.Println("Expected time to unlock : ", puzzle.TimeToUnlock)

	x       := base
	curr    := big.NewInt(0)
	one     := big.NewInt(1)
	zero    := big.NewInt(0)
	hundred := big.NewInt(100)
	start   := time.Now()

	for curr.Cmp(iter) < 0 {
		if mod(curr, hundred).Cmp(zero) == 0 {
			diff := big.NewInt(0).Sub(iter, curr)
			rem := time.Duration(int(diff.Int64())) * time.Second
			fmt.Println("Iteration", curr.String()+"/"+iter.String(), ", remaining time is", rem)
		}
		
		x = ExpByPowOfTwoModular(x, exp, n)

		curr.Add(curr,one)
	}

	fmt.Println("Puzzle solved, verifying...")
	sha := sha256AndHex(x.Bytes())
	shaOfSha := sha256AndHex([]byte(sha))

	if shaOfSha != puzzle.Sha256OfSha256 {
		fmt.Println("An error occured. The solution isn't the expected one.")
	} else {
		timeNeeded := time.Now().Sub(start)
		fmt.Println("################")
		fmt.Println("Time-Lock Puzzle solved in", timeNeeded, "!")
		fmt.Println("")
		fmt.Println(sha)
		fmt.Println("")
		fmt.Println("################")
	}
}