package main

import (
	"coin-tool/block"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("method", "", "method: getBalance, sendToken, sendEthWithSameAmount, sendEthWithDiffAmount, sendEthWithMultiToOne, sendEthMultiToMulti")
	pflag.String("rpc-url", "http://127.0.0.1:8545", "ethereum url to connect")
	pflag.StringSlice("from-address", []string{}, "transaction send address list")
	pflag.StringSlice("from-key", []string{}, "transaction send private key list")
	pflag.StringSlice("to-address", []string{}, "transaction to address list")
	pflag.StringSlice("diffAmount", []string{}, "transaction amount")
	pflag.String("contract", "0x2aC3c1d3e24b45c6C310534Bc2Dd84B5ed576335", "ethereum contract address")
	pflag.Int("decimals", 18, "ethereum account decimals")
	pflag.Float64("amount", 0, "ethereum account number")
	pflag.Int64("chain-id", 5, "ethereum chain id")
	pflag.Int64("gas-limit", 90000, "ethereum account number gas limit")
	pflag.Int64("gas-price", 50000000000, "ethereum account number gas price")

	_ = viper.BindPFlags(pflag.CommandLine)
}

func main() {
	pflag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed:%v", err)
	}

	params := block.CommandParams{
		Method:      viper.GetString("method"),
		RpcUrl:      viper.GetString("rpc-url"),
		FromAddress: viper.GetStringSlice("from-address"),
		FromKey:     viper.GetStringSlice("from-key"),
		ToAddress:   viper.GetStringSlice("to-address"),
		DiffAmount:  viper.GetStringSlice("diffAmount"),
		Contract:    viper.GetString("contract"),
		Amount:      viper.GetFloat64("amount"),
		ChainID:     viper.GetInt64("chain-id"),
		Decimals:    viper.GetInt("decimals"),
		GasLimit:    viper.GetInt64("gas-limit"),
		GasPrice:    viper.GetInt64("gas-price"),
	}

	log.Printf("command params:%+v\n\n", params)

	switch params.Method {
	case "sendToken", "sendEthWithSameAmount", "sendEthWithMultiToOne", "sendEthWithDiffAmount", "sendEthMultiToMulti":
		if len(params.FromAddress) == 0 || len(params.FromKey) == 0 || len(params.ToAddress) == 0 || params.Amount == 0 {
			log.Fatal("fromAddress, fromKey, toAddress, amount is required")
		}
		if len(params.FromAddress) != len(params.FromKey) {
			log.Fatal("from address and from key do not match")
		}

		sendTransaction(&params)
	case "getBalance":
		block.GetBalance(&params)
	default:
		log.Println("please input the correct method")
	}

	log.Println(params.Method, "execute finished")
}

func sendTransaction(params *block.CommandParams) {

	if params.Method == "sendEthWithSameAmount" {
		block.SendEthWithSameAmount(params)
		return
	}

	if params.Method == "sendEthWithDiffAmount" {
		block.SendEthWithDiffAmount(params)
		return
	}

	if params.Method == "sendEthWithMultiToOne" { //归集
		block.SendEthWithMultiToOne(params)
		return
	}

	if params.Method == "sendEthMultiToMulti" { //多对多发送余额
		block.SendEthWithMultiToMulti(params)
		return
	}

	if len(params.FromAddress) == 1 {
		block.SendTokenWithAmount(params)
		return
	} else if len(params.ToAddress) == 1 {
		block.SendTokenWithBalance(params)
		return
	}
}
