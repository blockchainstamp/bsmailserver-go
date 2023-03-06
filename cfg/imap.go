package cfg

import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"path/filepath"
)

type IMAPCfg struct {
	TlsKey         string      `json:"tls_key"`
	TlsCert        string      `json:"tls_cert"`
	SrvPort        int         `json:"srv_port"`
	SrvAddr        string      `json:"srv_addr"`
	AllowNotSecure bool        `json:"allow_not_secure"`
	TlsCfg         *tls.Config `json:"-"`
}

func (c *IMAPCfg) prepare(cfg, fPath string) error {
	if !filepath.IsAbs(fPath) {
		fPath = filepath.Join(cfg, string(filepath.Separator), fPath)
	}
	if err := util.ReadJsonFile(fPath, c); err != nil {
		fmt.Println("parse imap config failed:=>", err)
		return err
	}
	if c.AllowNotSecure {
		return nil
	}
	var (
		cPath = ""
		kPath = ""
	)
	if !filepath.IsAbs(c.TlsCert) {
		cPath = filepath.Join(cfg, string(filepath.Separator), c.TlsCert)
	}
	if !filepath.IsAbs(c.TlsKey) {
		kPath = filepath.Join(cfg, string(filepath.Separator), c.TlsKey)
	}
	tlsCfg, err := util.LoadServerTlsCnf(cPath, kPath)
	if err != nil {
		fmt.Println("load tls config of imap server failed:", err)
		return err
	}
	c.TlsCfg = tlsCfg
	return err
}