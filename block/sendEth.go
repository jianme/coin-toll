package block

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func SendEthWithSameAmount(params *CommandParams) {
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

	balance, err := ToWei(fmt.Sprintf("%v", params.Amount), ETH_DECIMALS)
	if err != nil {
		log.Fatalf("ToWei err:%v", err)
	}

	nonce, err := getTransactionCount(params.RpcUrl, params.FromAddress[0])
	if err != nil {
		log.Fatalf("getTransactionCount err:%v", err)
	}

	for _, toAddress := range params.ToAddress {
		tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), balance, uint64(params.GasLimit), big.NewInt(params.GasPrice), nil)
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
		time.Sleep(2 * time.Second)
	}
}

func SendEthWithDiffAmount(params *CommandParams) {
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

	nonce, err := getTransactionCount(params.RpcUrl, params.FromAddress[0])
	if err != nil {
		log.Fatalf("getTransactionCount err:%v", err)
	}

	for cnt, toAddress := range params.ToAddress {

		balance, err := ToWei(fmt.Sprintf("%v", params.DiffAmount[cnt]), ETH_DECIMALS)
		if err != nil {
			log.Fatalf("ToWei err:%v", err)
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), balance, uint64(params.GasLimit), big.NewInt(params.GasPrice), nil)
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
		time.Sleep(2 * time.Second)
	}
}

func SendEthWithMultiToMulti(params *CommandParams) {
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
		balance := big.NewInt(0).Sub(bigAmount, bigFee)
		if balance.Cmp(big.NewInt(0)) == -1 {
			continue
			log.Fatalf("insufficient balance:%v", balance)
		}

		nonce, err := getTransactionCount(params.RpcUrl, fromAddress)
		if err != nil {
			log.Fatalf("getTransactionCount err:%v", err)
			return
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(params.ToAddress[index]), balance, uint64(params.GasLimit), big.NewInt(params.GasPrice), nil)
		raw, err := SignTransaction(params.ChainID, tx, params.FromKey[index])
		if err != nil {
			log.Fatalf("SignTransaction err:%v", err)
		}

		txID, err := ethSendRawTransaction(params.RpcUrl, raw)
		if err != nil {
			log.Fatalf("ethSendRawTransaction err:%v", err)
		}

		log.Println(fromAddress, params.ToAddress[index], txID)
		//time.Sleep(time.Second)
	}
}

func SendEthWithMultiToOne(params *CommandParams) {
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
		balance := big.NewInt(0).Sub(bigAmount, bigFee)
		if balance.Cmp(big.NewInt(0)) == -1 {
			log.Fatalf("insufficient balance:%v", balance)
		}

		nonce, err := getTransactionCount(params.RpcUrl, fromAddress)
		if err != nil {
			log.Fatalf("getTransactionCount err:%v", err)
			return
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(params.ToAddress[0]), balance, uint64(params.GasLimit), big.NewInt(params.GasPrice), nil)
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

func getTransactionCount(rpcUrl, address string) (nonce uint64, err error) {
	jsonStr := fmt.Sprintf(`{"jsonrpc": "2.0", "id":1, "method": "eth_getTransactionCount", "params": ["%v", "pending"]}`, address)
	payload := strings.NewReader(jsonStr)
	client := &http.Client{}

	req, err := http.NewRequest("POST", rpcUrl, payload)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var resp BalanceResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	if len(resp.Result) <= 2 {
		err = fmt.Errorf("result length error")
		return
	}

	bigCount, boole := big.NewInt(0).SetString(resp.Result[2:], 16)
	if boole == false {
		err = fmt.Errorf("NewInt SetString err")
		return
	}

	return bigCount.Uint64(), nil
}

func ethSendRawTransaction(rpcUrl, raw string) (txID string, err error) {
	jsonStr := fmt.Sprintf(`{"jsonrpc": "2.0", "id":1, "method": "eth_sendRawTransaction", "params": ["%v"]}`, raw)
	payload := strings.NewReader(jsonStr)
	client := &http.Client{}

	req, err := http.NewRequest("POST", rpcUrl, payload)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var resp BalanceResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	return resp.Result, nil
}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignTransaction(chainID int64, tx *types.Transaction, privateKeyStr string) (string, error) {
	privateKey, err := StringToPrivateKey(fmt.Sprintf("0x%v", privateKeyStr))
	if err != nil {
		return "", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%v", hex.EncodeToString(b)), nil
}
