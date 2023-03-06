package service

type SimpleMailSrv struct {
	smtpSrv *SmtpSrv
	imapSrv *ImapSrv
}

func (ss *SimpleMailSrv) Run() {
	ss.smtpSrv = newSmtpSrv()
	ss.imapSrv = newImapSrv()
	ss.smtpSrv.Run()
	ss.imapSrv.Run()
	_srvLog.Info("simple mail service start success......")
}

func (ss *SimpleMailSrv) ShutDown() {
	if ss.smtpSrv != nil {
		_ = ss.smtpSrv.Close()
	}
	if ss.imapSrv != nil {
		_ = ss.imapSrv.Close()
	}
	_srvLog.Info("simple mail service shutting down......")
}
