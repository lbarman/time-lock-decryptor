package main

import (
    "math/big"
    "crypto/sha256"
    "encoding/hex"
)

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