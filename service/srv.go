package service

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	_inst   Service = nil
	_once   sync.Once
	_srvLog = logrus.WithFields(logrus.Fields{
		"mode": "mail service",
	})
)

type Service interface {
	Run()
	ShutDown()
}

func Inst() Service {
	_once.Do(func() {
		if cfg.UseStamp() {
			_inst = bsMailSrv()
		} else {
			_inst = newSimpleMail()
		}
	})
	return _inst
}

func newSimpleMail() Service {
	fmt.Println("using simple mail service......")
	sm := &SimpleMailSrv{}
	return sm
}

func bsMailSrv() Service {
	fmt.Println("using block chain stamp......")
	bsm := &BlockchainMailSrv{}
	return bsm
}
