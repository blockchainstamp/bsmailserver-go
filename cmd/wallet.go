package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "wallet action",
		Long:  wUsage,
		Run:   walletAct,
	}
	pwd    = ""
	wUsage = `
	usage:
		wallet create --pwd create new wallet with password [pwd]
		wallet show --list show current wallet
`
)

func init() {

	walletCmd.Flags().StringVarP(&pwd, "pwd", "p", "",
		"password of wallet")

	rootCmd.AddCommand(walletCmd)
}

func walletAct(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		fmt.Println(cmd.Long)
		return
	}
	var actTyp = args[0]
	switch actTyp {
	case "create":
		fmt.Println("start to create wallet")
		break
	case "show":
		fmt.Println("start to show wallet")
		break
	default:
		fmt.Println(cmd.Long)
	}
}
