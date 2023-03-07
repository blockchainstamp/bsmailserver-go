package cfg

import (
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"os"
	"path/filepath"
)

type SMTPCfg struct {
	SrvAddr         string      `json:"srv_addr"`
	SrvDomain       string      `json:"srv_domain"`
	TlsKey          string      `json:"tls_key"`
	TlsCert         string      `json:"tls_cert"`
	DKIMKey         string      `json:"dkim_key"`
	SrvPort         int         `json:"srv_port"`
	AllowNotSecure  bool        `json:"allow_not_secure"`
	MaxMessageBytes int         `json:"max_message_bytes"`
	ReadTimeOut     int         `json:"read_time_out"`
	WriteTimeOut    int         `json:"write_time_out"`
	MaxRecipients   int         `json:"max_recipients"`
	TlsCfg          *tls.Config `json:"-"`
	DKIMKeyData     []byte      `json:"-"`
}

func (c *SMTPCfg) prepare(homeDir, fPath string) error {
	if !filepath.IsAbs(fPath) {
		fPath = filepath.Join(homeDir, string(filepath.Separator), fPath)
	}
	if err := util.ReadJsonFile(fPath, c); err != nil {
		fmt.Println("parse smtp config failed:=>", err)
		return err
	}
	if c.AllowNotSecure {
		return nil
	}
	var (
		cPath = ""
		kPath = ""
		dPath = ""
	)
	if !filepath.IsAbs(c.TlsCert) {
		cPath = filepath.Join(homeDir, string(filepath.Separator), c.TlsCert)
	}
	if !filepath.IsAbs(c.TlsKey) {
		kPath = filepath.Join(homeDir, string(filepath.Separator), c.TlsKey)
	}
	tlsCfg, err := util.LoadServerTlsCnf(cPath, kPath)
	if err != nil {
		fmt.Println("load tls config of smtp server failed:", err)
		return err
	}
	c.TlsCfg = tlsCfg

	if !filepath.IsAbs(c.DKIMKey) {
		dPath = filepath.Join(homeDir, string(filepath.Separator), c.DKIMKey)
	}
	bts, err := os.ReadFile(dPath)
	if err != nil {
		fmt.Println("load DKIM key for smtp server failed:", err)
		return err
	}
	c.DKIMKeyData = bts
	return err
}
