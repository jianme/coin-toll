package block

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type CommandParams struct {
	Method      string
	RpcUrl      string
	FromAddress []string
	FromKey     []string
	ToAddress   []string
	DiffAmount  []string
	Contract    string
	Decimals    int
	Amount      float64
	ChainID     int64
	GasLimit    int64
	GasPrice    int64
}

var ETH_DECIMALS = 18
var ERC20_DECIMALS = 8

func GetBigFromHex(hexamount string) (*big.Int, error) {
	s := hexamount
	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	if len(s) == 0 {
		return big.NewInt(0), nil
	}

	n, boole := big.NewInt(0).SetString(s, 16)
	if boole == false {
		err := errors.New("invalid syntax")
		return nil, err
	}

	return n, nil
}

// wei -> count
func EthFromWei(bigs *big.Int) (strvalue string, err error) {
	strvalue = "0"
	decimals := ETH_DECIMALS

	strn := bigs.String()
	if strn == "0" {
		return
	}

	if len(strn) > decimals {
		if decimals > ERC20_DECIMALS {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+ERC20_DECIMALS])
		} else {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+decimals])
		}
	} else {
		if decimals > ERC20_DECIMALS {
			if len(strn) > decimals-ERC20_DECIMALS {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", ERC20_DECIMALS, "s")
				strvalue = fmt.Sprintf(strFormat, strn[:len(strn)-(decimals-ERC20_DECIMALS)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strvalue = fmt.Sprintf(strFormat, strn)
		}
	}

	return
}

// wei -> count
func FromWeiWithDecimals(s string, decimals int) (strvalue string, err error) {

	strvalue = "0"

	if len(s) < 1 {
		err = errors.New("invalid syntax")
		return
	}

	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	if len(s) == 0 {
		return
	}

	n, boole := big.NewInt(0).SetString(s, 16)
	if boole == false {
		err = errors.New("invalid syntax")
		return
	}

	strn := n.String()
	if strn == "0" {
		return
	}

	if len(strn) > decimals {
		if decimals > ERC20_DECIMALS {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+ERC20_DECIMALS])
		} else {
			strvalue = fmt.Sprintf("%v.%v", strn[:len(strn)-decimals], strn[len(strn)-decimals:len(strn)-decimals+decimals])
		}
	} else {
		if decimals > ERC20_DECIMALS {
			if len(strn) > decimals-ERC20_DECIMALS {
				strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", ERC20_DECIMALS, "s")
				strvalue = fmt.Sprintf(strFormat, strn[:len(strn)-(decimals-ERC20_DECIMALS)])
			}
		} else {
			strFormat := fmt.Sprintf(`0.%v%v%v`, "%0", decimals, "s")
			strvalue = fmt.Sprintf(strFormat, strn)
		}
	}

	return
}

func ToWei(svalue string, decimals int) (amount *big.Int, err error) {
	bigdecimals := big.NewInt(10)
	bigdecimals = bigdecimals.Exp(bigdecimals, big.NewInt(int64(decimals)), nil)

	priceparts := strings.Split(svalue, ".")
	if len(priceparts) == 1 {
		count, err1 := strconv.ParseUint(svalue, 10, 64)
		if err != nil {
			err = err1
			return
		}

		iamount := big.NewInt(int64(count))
		amount = iamount.Mul(iamount, bigdecimals)
		return
	} else if len(priceparts) == 2 {
		var uprice1 uint64
		var uprice2 uint64
		uprice1, err1 := strconv.ParseUint(priceparts[0], 10, 64)
		if err1 != nil {
			err = err1
			return
		}

		if len(priceparts[1]) <= 0 {
			iamount := big.NewInt(int64(uprice1))
			amount = iamount.Mul(iamount, bigdecimals)
			return
		}

		if len(priceparts[1]) > decimals {
			priceparts[1] = priceparts[1][:decimals]
		}
		uprice2, err = strconv.ParseUint(priceparts[1], 10, 64)
		if err != nil {
			return
		}

		iamount1 := big.NewInt(int64(uprice1))
		iamount1.Mul(iamount1, bigdecimals)
		iamount2 := big.NewInt(int64(uprice2))
		iamount2.Mul(iamount2, bigdecimals.Exp(big.NewInt(10), big.NewInt(int64(decimals-len(priceparts[1]))), nil))
		amount = iamount1.Add(iamount1, iamount2)
		return
	} else {
		err = errors.New("input invalid")
		return
	}

	return
}
