package cfg

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"path/filepath"
)

type SMTPCfg struct {
	SrvAddr         string        `json:"srv_addr"`
	SrvDomain       string        `json:"srv_domain"`
	TlsKey          string        `json:"tls_key"`
	TlsCert         string        `json:"tls_cert"`
	DKIMKey         string        `json:"dkim_key"`
	DKIMSelector    string        `json:"dkim_selector"`
	SrvPort         int           `json:"srv_port"`
	AllowNotSecure  bool          `json:"allow_not_secure"`
	MaxMessageBytes int           `json:"max_message_bytes"`
	ReadTimeOut     int           `json:"read_time_out"`
	WriteTimeOut    int           `json:"write_time_out"`
	MaxRecipients   int           `json:"max_recipients"`
	TlsCfg          *tls.Config   `json:"-"`
	DKIMSigner      crypto.Signer `json:"-"`
}

func (c *SMTPCfg) String() string {
	s := "\n======SMTP Config======"
	s += "\nSrvAddr:  \t" + c.SrvAddr
	s += "\nSrvDomain:\t" + c.SrvDomain
	s += "\nTlsKey:   \t" + c.TlsKey
	s += "\nTlsCert:\t" + c.TlsCert
	s += "\nDKIMKey:\t" + c.DKIMKey
	s += fmt.Sprintf("\nSrvPort:\t%d", c.SrvPort)
	s += fmt.Sprintf("\nAllowNotSecure:\t%t", c.AllowNotSecure)
	s += fmt.Sprintf("\nMaxMessageBytes:\t%d", c.MaxMessageBytes)
	s += fmt.Sprintf("\nReadTimeOut:\t%d", c.ReadTimeOut)
	s += fmt.Sprintf("\nWriteTimeOut:\t%d", c.WriteTimeOut)
	s += fmt.Sprintf("\nMaxRecipients:\t%d", c.MaxRecipients)
	s += "\n========================="
	return s
}

func (c *SMTPCfg) prepare(homeDir, fPath string) error {
	if !filepath.IsAbs(fPath) {
		fPath = filepath.Join(homeDir, string(filepath.Separator), fPath)
	}
	if err := util.ReadJsonFile(fPath, c); err != nil {
		fmt.Println("parse smtp config failed:=>", err)
		return err
	}
	fmt.Println(c.String())

	var (
		cPath = c.TlsCert
		kPath = c.TlsKey
		dPath = c.DKIMKey
	)
	if !c.AllowNotSecure && (c.TlsCert == "" || c.TlsKey == "") {
		return util.CfgNoTlsFile
	}

	if c.TlsCert != "" && c.TlsKey != "" {
		if !filepath.IsAbs(c.TlsCert) {
			cPath = filepath.Join(homeDir, string(filepath.Separator), c.TlsCert)
		}
		if !filepath.IsAbs(c.TlsKey) {
			kPath = filepath.Join(homeDir, string(filepath.Separator), c.TlsKey)
		}
		tlsCfg, err := util.LoadServerTlsCnf(cPath, kPath)
		if err != nil {
			fmt.Println("load tls config of smtp server failed:", cPath, kPath, err)
			return err
		}
		c.TlsCfg = tlsCfg
	}

	if !filepath.IsAbs(c.DKIMKey) {
		dPath = filepath.Join(homeDir, string(filepath.Separator), c.DKIMKey)
	}

	signer, err := util.LoadDkimKey(dPath)
	if err != nil {
		return err
	}
	c.DKIMSigner = signer
	return nil
}
