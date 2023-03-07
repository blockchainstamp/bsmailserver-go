package cfg

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	_curSysConf = &SysRunTimeConfig{}
)

type SysStaticConfig struct {
	SmtpCfg     string `json:"smtp_cfg"`
	ImapCfg     string `json:"imap_cfg"`
	BSCfg       string `json:"bs_cfg"`
	BackendCfg  string `json:"backend_cfg"`
	WalletInUse string `json:"wallet_in_use"`
	LogLevel    string `json:"log_level"`
	DBPath      string `json:"db_path"`
	CmdSrvAddr  string `json:"cmd_srv_addr"`
}

type SysRunTimeConfig struct {
	useStamp bool
	smtp     *SMTPCfg
	imap     *IMAPCfg
	bStamp   *BStampConf
	backend  *BackConfig
}

func PrepareConfig(homeDir string) error {
	var (
		err         error
		sc          = &SMTPCfg{}
		ic          = &IMAPCfg{}
		bc          = &BStampConf{}
		backCfg     = &BackConfig{}
		sysConfPath = ""
	)

	if homeDir == "" {
		homeDir = util.DefaultBaseDir
	}
	sysConfPath = filepath.Join(homeDir, string(filepath.Separator), util.DefaultSysConfig)
	c := &SysStaticConfig{}
	if err = util.ReadJsonFile(sysConfPath, c); err != nil {
		fmt.Println("parse system config failed:=>", err)
		return err
	}
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		fmt.Println("set system log level failed:=>", err)
		return err
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err = sc.prepare(homeDir, c.SmtpCfg)
	if err != nil {
		return err
	}
	err = ic.prepare(homeDir, c.ImapCfg)
	if err != nil {
		return err
	}

	err = bc.prepare(homeDir, c.BSCfg)
	if err != nil {
		return err
	}

	err = backCfg.prepare(homeDir, c.BackendCfg)
	if err != nil {
		return err
	}

	_curSysConf = &SysRunTimeConfig{
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

func CurSrvConf() *SysRunTimeConfig {
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

func CurBackendConf() *BackConfig {
	if _curSysConf.backend == nil {
		panic("init stamp config first please!!!")
	}
	return _curSysConf.backend
}
