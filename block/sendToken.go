package block

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
)

func SendTokenWithAmount(params *CommandParams) {
	bigFee := big.NewInt(0).Mul(big.NewInt(params.GasLimit), big.NewInt(params.GasPrice))
	bigFee = big.NewInt(0).Mul(bigFee, big.NewInt(int64(len(params.ToAddress))))

	amount, err := getBalance(params.RpcUrl, params.FromAddress[0])
	if err != nil {
		log.Fatalf("getBalance err:%v", err)
	}

	bigAmount, err := GetBigFromHex(amount)
	if err != nil {
		log.Fatalf("GetBigFromHex err:%v", err)
	}

	subFeeAmount := big.NewInt(0).Sub(bigAmount, bigFee)
	if subFeeAmount.Cmp(big.NewInt(0)) == -1 {
		log.Fatalf("insufficient balance:%v", subFeeAmount)
	}

	balance, err := ToWei(fmt.Sprintf("%v", params.Amount), params.Decimals)
	if err != nil {
		log.Fatalf("ToWei err:%v", err)
	}

	nonce, err := getTransactionCount(params.RpcUrl, params.FromAddress[0])
	if err != nil {
		log.Fatalf("getTransactionCount err:%v", err)
	}

	for _, toAddress := range params.ToAddress {
		data, err := MakeERC20TransferData(toAddress, balance)
		if err != nil {
			log.Fatalf("MakeERC20TransferData err:%v", err)
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(params.Contract), big.NewInt(0), uint64(params.GasLimit), big.NewInt(params.GasPrice), data)
		raw, err := SignTransaction(params.ChainID, tx, params.FromKey[0])
		if err != nil {
			log.Fatalf("SignTransaction err:%v", err)
		}

		txID, err := ethSendRawTransaction(params.RpcUrl, raw)
		if err != nil {
			log.Fatalf("ethSendRawTransaction err:%v", err)
		}

		log.Println(params.FromAddress[0], toAddress, txID)
		nonce++
	}
}

func SendTokenWithBalance(params *CommandParams) {
	for index, fromAddress := range params.FromAddress {
		amount, err := getBalance(params.RpcUrl, fromAddress)
		if err != nil {
			log.Fatalf("getBalance err:%v", err)
		}

		bigAmount, err := GetBigFromHex(amount)
		if err != nil {
			log.Fatalf("GetBigFromHex err:%v", err)
		}

		bigFee := big.NewInt(0).Mul(big.NewInt(params.GasLimit), big.NewInt(params.GasPrice))
		subFeeAmount := big.NewInt(0).Sub(bigAmount, bigFee)
		if subFeeAmount.Cmp(big.NewInt(0)) == -1 {
			log.Fatalf("insufficient balance:%v", subFeeAmount)
		}

		tokenAmount, err := getTokenBalance(params.RpcUrl, fromAddress, params.Contract)
		if err != nil {
			log.Fatalf("getTokenBalance err:%v", err)
		}

		balance, err := GetBigFromHex(tokenAmount)
		if err != nil {
			log.Fatalf("GetBigFromHex err:%v", err)
		}

		nonce, err := getTransactionCount(params.RpcUrl, fromAddress)
		if err != nil {
			log.Fatalf("getTransactionCount err:%v", err)
			return
		}

		data, err := MakeERC20TransferData(params.ToAddress[0], balance)
		if err != nil {
			log.Fatalf("MakeERC20TransferData err:%v", err)
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(params.Contract), big.NewInt(0), uint64(params.GasLimit), big.NewInt(params.GasPrice), data)
		raw, err := SignTransaction(params.ChainID, tx, params.FromKey[index])
		if err != nil {
			log.Fatalf("SignTransaction err:%v", err)
		}

		txID, err := ethSendRawTransaction(params.RpcUrl, raw)
		if err != nil {
			log.Fatalf("ethSendRawTransaction err:%v", err)
		}

		log.Println(fromAddress, params.ToAddress[0], txID)
	}
}

func MakeERC20TransferData(toAddress string, amount *big.Int) ([]byte, error) {
	methodId := crypto.Keccak256Hash([]byte("transfer(address,uint256)"))
	var data []byte
	data = append(data, methodId[:4]...)
	paddedAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...)
	return data, nil
}
