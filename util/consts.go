package util

const (
	DefaultBaseDir    = "MailBase"
	DefaultSysConfig  = "bsmail.conf"
	DefaultConfDir    = "config"
	DefaultWalletDir  = "wallets"
	DefaultCmdSrvAddr = "127.0.0.1:10001"
	DefaultFilePerm   = 0777
	DefaultSMTPPort   = 465
	DefaultIMAPPort   = 993
	MaxMailSize       = 1 << 26
	SMTPReadTimeOut   = 10
	SMTPWriteTimeOut  = 10
	SMTPMaxRecipients = 1 << 10

	DBMemHome    = "memDB"
	DBJsonHome   = "jsonDB"
	DBLevelHome  = "levelDB"
	DBSqliteHome = "sqliteDB"
)
