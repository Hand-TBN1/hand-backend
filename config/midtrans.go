package config

import (
	"github.com/veritrans/go-midtrans"
)

func SetupMidtrans() midtrans.Client {
	client := midtrans.NewClient()
	client.ServerKey = Env.MidtransServerKey
	client.ClientKey = Env.MidtransClientKey
	client.APIEnvType = midtrans.Sandbox  // Or midtrans.Production based on the environment

	return client
}
