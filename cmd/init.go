package main

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/backstore"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "init system",
		Long:  `init system setting by default value`,
		Run:   initSystemDefault,
	}
	baseDir = ""
	force   bool
)

func init() {
	initCmd.Flags().StringVarP(&baseDir, "base", "b", util.DefaultBaseDir,
		"password of wallet")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "--force overwrite current data")
	rootCmd.AddCommand(initCmd)
}

func initSystemDefault(cmd *cobra.Command, args []string) {
	dc := cfg.SysStaticConfig{
		SmtpCfg:     "config/smtp.json",
		ImapCfg:     "config/imap.json",
		BSCfg:       "config/stamp.json",
		BackendCfg:  "config/backend.json",
		WalletInUse: "wallets/main.json",
		LogLevel:    "info",
		CmdSrvAddr:  util.DefaultCmdSrvAddr,
	}
	if !force {
		if _, has := util.FileExists(baseDir); has {
			fmt.Println("Duplicate Init!!! \nuse --force if you want to overwrite current base directory[Dangerous!!!]")
			return
		}
	}

	if err := util.TouchDir(baseDir); err != nil {
		panic(err)
	}
	confPath := path.Join(baseDir, string(filepath.Separator), util.DefaultConfDir)
	if err := util.TouchDir(confPath); err != nil {
		panic(err)
	}

	smtpCfgFilePath := path.Join(baseDir, string(filepath.Separator), dc.SmtpCfg)
	sConf := cfg.SMTPCfg{
		SrvAddr:         "0.0.0.0",
		SrvDomain:       "localhost", //smtp.simplenets.org
		TlsKey:          "config/key.pem",
		TlsCert:         "config/cert.pem",
		DKIMKey:         "config/dkim.pem",
		SrvPort:         util.DefaultSMTPPort,
		MaxMessageBytes: util.MaxMailSize,
		ReadTimeOut:     util.SMTPReadTimeOut,
		WriteTimeOut:    util.SMTPWriteTimeOut,
		MaxRecipients:   util.SMTPMaxRecipients,
	}

	if err := util.WriteJsonFile(smtpCfgFilePath, sConf); err != nil {
		panic(err)
	}

	imapCfgFilePath := path.Join(baseDir, string(filepath.Separator), dc.ImapCfg)
	iConf := cfg.IMAPCfg{
		TlsKey:  "config/key.pem",
		TlsCert: "config/cert.pem",
		SrvPort: util.DefaultIMAPPort,
		SrvAddr: "0.0.0.0",
	}
	if err := util.WriteJsonFile(imapCfgFilePath, iConf); err != nil {
		panic(err)
	}

	bsCfgFilePath := path.Join(baseDir, string(filepath.Separator), dc.BSCfg)
	bsConf := cfg.BStampConf{
		WalletPwd: "123",
	}
	if err := util.WriteJsonFile(bsCfgFilePath, bsConf); err != nil {
		panic(err)
	}

	backendCfg := cfg.BackConfig{
		DBType: cfg.DBTypMem,
		DBHome: "mail_data",
	}
	backedCfgFilePath := path.Join(baseDir, string(filepath.Separator), dc.BackendCfg)
	if err := util.WriteJsonFile(backedCfgFilePath, backendCfg); err != nil {
		panic(err)
	}
	dbHome := path.Join(baseDir, string(filepath.Separator), backendCfg.DBHome)
	if err := backstore.InitDefaultDB(dbHome); err != nil {
		panic(err)
	}

	wPath := path.Join(baseDir, string(filepath.Separator), util.DefaultWalletDir)
	if err := util.TouchDir(wPath); err != nil {
		panic(err)
	}

	sysCfgFilePath := path.Join(baseDir, string(filepath.Separator), util.DefaultSysConfig)
	if err := util.WriteJsonFile(sysCfgFilePath, dc); err != nil {
		panic(err)
	}
}
