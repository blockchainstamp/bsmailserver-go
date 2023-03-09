package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/backstore"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/util"
	iSrv "github.com/emersion/go-imap/server"
)

type ImapSrv struct {
	imapSrv    *iSrv.Server
	imap143Srv *iSrv.Server
}

func newImapSrv() *ImapSrv {
	var is = &ImapSrv{}
	var conf = cfg.CurImapConf()
	imapSrv := iSrv.New(backstore.Inst())
	imapSrv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, conf.SrvPort)
	imapSrv.AllowInsecureAuth = conf.AllowNotSecure
	imapSrv.TLSConfig = conf.TlsCfg
	is.imapSrv = imapSrv
	_srvLog.Info("imap service init success at: ", imapSrv.Addr)

	imap143Srv := iSrv.New(backstore.Inst())
	imap143Srv.Addr = fmt.Sprintf("%s:%d", conf.SrvAddr, util.DefaultSystemImapPort)
	imap143Srv.AllowInsecureAuth = conf.AllowNotSecure
	imap143Srv.TLSConfig = conf.TlsCfg
	is.imap143Srv = imap143Srv
	_srvLog.Info("default system imap service init success at: ", imap143Srv.Addr)

	return is
}

func (is *ImapSrv) Close() error {
	err1, err2 := is.imapSrv.Close(), is.imap143Srv.Close()
	if err2 == nil && err1 == nil {
		return nil
	}
	return fmt.Errorf("%v:%v", err1, err2)
}

func (is *ImapSrv) Run() {
	go func() {
		_srvLog.Info("starting imap service.......")
		err := is.imapSrv.ListenAndServeTLS()
		if err != nil {
			_srvLog.Error("imap service exit:", err)
			panic(err)
		}
	}()
	go func() {
		_srvLog.Info("starting system 143 imap service.......")
		err := is.imap143Srv.ListenAndServeTLS()
		if err != nil {
			_srvLog.Error("system 143 imap service exit:", err)
			panic(err)
		}
	}()
}
