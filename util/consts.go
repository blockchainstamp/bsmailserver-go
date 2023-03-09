package util

import "fmt"

const (
	DefaultBaseDir        = "MailBase"
	DefaultSysConfig      = "bsmail.conf"
	DefaultConfDir        = "config"
	DefaultWalletDir      = "wallets"
	DefaultCmdSrvAddr     = "127.0.0.1:10001"
	DefaultFilePerm       = 0777
	DefaultSMTPPort       = 465
	DefaultSystemSmtpPort = 25
	DefaultIMAPPort       = 993
	DefaultSystemImapPort = 143
	MaxMailSize           = 1 << 26
	SMTPReadTimeOut       = 10
	SMTPWriteTimeOut      = 10
	SMTPMaxRecipients     = 1 << 10

	DBMemHome    = "memDB"
	DBJsonHome   = "jsonDB"
	DBLevelHome  = "levelDB"
	DBSqliteHome = "sqliteDB"

	MHKeyDkimSign  = "Dkim-Signature"
	MHKeyStampSign = "X-Stamp"
	MHKeyMsgID     = "Message-Id"
	MHKeyFrom      = "From"
	MHKeyTos       = "To"

	MBXIndex = "INBOX"
)

var (
	BackendNoSuchUser       = fmt.Errorf("no such user")
	BackendBadUser          = fmt.Errorf("bad user name and password")
	BackendNoSuchMailBox    = fmt.Errorf("no such mail box")
	BackendDuplicateMailBox = fmt.Errorf("duplicated mail box")
	BackendDeleteInbox      = fmt.Errorf("inbox can't be removed")
	CfgNoTlsFile            = fmt.Errorf("invalid tls file config")
	CfgInvalidDkim          = fmt.Errorf("invalid dkim file")
	SMTPNoValidMX           = fmt.Errorf("no valid mx record")
)
