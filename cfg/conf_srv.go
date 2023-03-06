package cfg

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	_curSysConf = &SrvRunTimeConfig{}
)

type SrvStaticConfig struct {
	SmtpCfg     string `json:"smtp_cfg"`
	ImapCfg     string `json:"imap_cfg"`
	BSCfg       string `json:"bs_cfg"`
	WalletInUse string `json:"wallet_in_use"`
	LogLevel    string `json:"log_level"`
	DBPath      string `json:"db_path"`
	CmdSrvAddr  string `json:"cmd_srv_addr"`
}

type SrvRunTimeConfig struct {
	useStamp bool
	smtp     *SMTPCfg
	imap     *IMAPCfg
	bStamp   *BStampConf
}

func PrepareConfig(dir string) error {
	var (
		err         error
		sysConfPath = ""
		confPath    = ""
	)

	if dir == "" {
		sysConfPath = filepath.Join(util.DefaultBaseDir, string(filepath.Separator), util.DefaultSysConfig)
	} else {
		sysConfPath = filepath.Join(dir, string(filepath.Separator), util.DefaultSysConfig)
	}

	c := &SrvStaticConfig{}
	if err = util.ReadJsonFile(sysConfPath, c); err != nil {
		fmt.Println("parse system config failed:=>", err)
		return err
	}
	confPath = filepath.Join(dir, string(filepath.Separator), util.DefaultConfDir)

	sc := &SMTPCfg{}
	err = sc.prepare(confPath, c.SmtpCfg)
	if err != nil {
		return err
	}
	ic := &IMAPCfg{}
	err = ic.prepare(confPath, c.ImapCfg)
	if err != nil {
		return err
	}

	bc := &BStampConf{}
	err = bc.prepare(confPath, c.BSCfg)
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		fmt.Println("set system log level failed:=>", err)
		return err
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	_curSysConf = &SrvRunTimeConfig{
		smtp:   sc,
		imap:   ic,
		bStamp: bc,
	}
	return nil
}
func UseStamp() bool {
	return _curSysConf.useStamp
}

func SetStampTag(sInUse bool) {
	_curSysConf.useStamp = sInUse
}

func CurSrvConf() *SrvRunTimeConfig {
	return _curSysConf
}

func CurSmtpConf() *SMTPCfg {
	if _curSysConf.smtp == nil {
		panic("init smtp config first please!!!")
	}
	return _curSysConf.smtp
}

func CurImapConf() *IMAPCfg {
	if _curSysConf.imap == nil {
		panic("init imap config first please!!!")
	}
	return _curSysConf.imap
}

func CurBStampConf() *BStampConf {
	if _curSysConf.bStamp == nil {
		panic("init stamp config first please!!!")
	}
	return _curSysConf.bStamp
}
