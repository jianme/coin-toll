package block

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type BalanceResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func GetBalance(params *CommandParams) {
	for _, address := range params.FromAddress {
		bigAmount, err := getBalance(params.RpcUrl, address)
		if err != nil {
			log.Fatalf("getBalance err:%v", err)
			return
		}

		balance, err := FromWeiWithDecimals(bigAmount, ETH_DECIMALS)
		if err != nil {
			log.Fatalf("FromWeiWithDecimals err:%v", err)
			return
		}
		fmt.Println(balance)
		/*
			bigAmount, err = getTokenBalance(params.RpcUrl, address, params.Contract)
			if err != nil {
				log.Fatalf("getTokenBalance err:%v", err)
				return
			}

			balance, err = FromWeiWithDecimals(bigAmount, params.Decimals)
			if err != nil {
				log.Fatalf("FromWeiWithDecimals err:%v", err)
				return
			}
			log.Println(address, balance)
		*/
	}
}

func getBalance(rpcUrl, address string) (amount string, err error) {
	jsonStr := fmt.Sprintf(`{"jsonrpc": "2.0", "id":1, "method": "eth_getBalance", "params": ["%v","latest"]}`, address)
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

func getTokenBalance(rpcUrl, address string, contract string) (amount string, err error) {
	jsonData := fmt.Sprintf("%v%064s", "0x70a08231", address[2:])
	jsonStr := fmt.Sprintf(`{"jsonrpc": "2.0", "id":1, "method": "eth_call", "params": [{"to":"%v", "data":"%v"},"latest"]}`, contract, jsonData)
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
