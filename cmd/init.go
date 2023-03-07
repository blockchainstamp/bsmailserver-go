package main

import (
	"fmt"
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
		DBPath:      "mail_data",
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

	subPath := path.Join(baseDir, string(filepath.Separator), dc.SmtpCfg)
	sConf := cfg.SMTPCfg{
		SrvAddr:         "0.0.0.0",
		SrvDomain:       "smtp.simplenets.org",
		TlsKey:          "/etc/letsencrypt/live/smtp.simplenets.org/privkey.pem",
		TlsCert:         "/etc/letsencrypt/live/smtp.simplenets.org/fullchain.pem",
		DKIMKey:         "dkim.key",
		SrvPort:         util.DefaultSMTPPort,
		MaxMessageBytes: util.MaxMailSize,
		ReadTimeOut:     util.SMTPReadTimeOut,
		WriteTimeOut:    util.SMTPWriteTimeOut,
		MaxRecipients:   util.SMTPMaxRecipients,
	}

	if err := util.WriteJsonFile(subPath, sConf); err != nil {
		panic(err)
	}

	subPath = path.Join(baseDir, string(filepath.Separator), dc.ImapCfg)
	iConf := cfg.IMAPCfg{
		TlsKey:  "/etc/letsencrypt/live/smtp.simplenets.org/privkey.pem",
		TlsCert: "/etc/letsencrypt/live/smtp.simplenets.org/fullchain.pem",
		SrvPort: util.DefaultIMAPPort,
		SrvAddr: "0.0.0.0",
	}
	if err := util.WriteJsonFile(subPath, iConf); err != nil {
		panic(err)
	}

	subPath = path.Join(baseDir, string(filepath.Separator), dc.BSCfg)
	bsConf := cfg.BStampConf{
		WalletPwd: "123",
	}
	if err := util.WriteJsonFile(subPath, bsConf); err != nil {
		panic(err)
	}

	subPath = path.Join(baseDir, string(filepath.Separator), dc.BackendCfg)
	backendCfg := cfg.BackConfig{
		UseMemDB: true,
	}
	if err := util.WriteJsonFile(subPath, backendCfg); err != nil {
		panic(subPath)
	}

	subPath = path.Join(baseDir, string(filepath.Separator), dc.DBPath)
	if err := util.TouchDir(subPath); err != nil {
		panic(err)
	}

	wPath := path.Join(baseDir, string(filepath.Separator), util.DefaultWalletDir)
	if err := util.TouchDir(wPath); err != nil {
		panic(err)
	}

	subPath = path.Join(baseDir, string(filepath.Separator), util.DefaultSysConfig)
	if err := util.WriteJsonFile(subPath, dc); err != nil {
		panic(err)
	}
}
