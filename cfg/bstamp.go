package cfg

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"path/filepath"
)

type BStampConf struct {
	WalletPwd string `json:"wallet_pwd"`
}

func (c *BStampConf) prepare(cfg, fPath string) error {
	if !filepath.IsAbs(fPath) {
		fPath = filepath.Join(cfg, string(filepath.Separator), fPath)
	}
	if err := util.ReadJsonFile(fPath, c); err != nil {
		fmt.Println("parse blockchain stamp config failed:=>", err)
		return err
	}
	fmt.Println(c.String())
	return nil
}

func (c *BStampConf) String() string {
	s := "\n======Blockchain Stamp Config======"
	s += "\nWalletPwd:\t***"
	s += "\n==================================="
	return s
}
