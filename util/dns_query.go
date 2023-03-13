package util

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

var (
	_inst   *DnsUtil = nil
	_once   sync.Once
	_dnsLog = logrus.WithFields(logrus.Fields{
		"mode": "dns query",
	})
)

type DnsUtil struct {
	sync.RWMutex
	MXs map[string][]*net.MX
}

func DnsInst() *DnsUtil {
	_once.Do(func() {
		_inst = &DnsUtil{MXs: make(map[string][]*net.MX)}
	})
	return _inst
}

func tryConnect(lDomain, rHost string, mx *net.MX) (c *smtp.Client, err error) {
	_dnsLog.Debugf("prepare to try mx:%+v", mx)
	addr := fmt.Sprintf("%s:%d", mx.Host, DefaultSystemSmtpPort)
	conn, err := net.DialTimeout("tcp", addr, MailMTATimeOut)
	if err != nil {
		_dnsLog.Warnf("dial(%s) err: %s", addr, err)
		return
	}
	c, err = smtp.NewClient(conn, mx.Host)
	if err != nil {
		_dnsLog.Warnf("dial(%s) err: %s", addr, err)
		goto closeAndRet
	}
	err = c.Hello(lDomain)
	if err != nil {
		_dnsLog.Warn("say hello err:", err, lDomain, mx.Host)
		goto closeAndRet
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: rHost}
		err = c.StartTLS(tlsCfg)
		if err != nil {
			return nil, err
		}
	}

	return c, nil

closeAndRet:
	conn.Close()
	return nil, err
}

func (du *DnsUtil) ValidSmtpCli(lDomain, rDomain string) (c *smtp.Client, err error) {
	var (
		mxs []*net.MX
		ok  bool
	)

	du.RLock()
	mxs, ok = du.MXs[rDomain]
	du.RUnlock()
	if !ok {
		_dnsLog.Info("no cached mx record for domain:", rDomain)
		mxs, err = net.LookupMX(rDomain)
		if err != nil {
			_dnsLog.Warn("LookupMX err:", err, rDomain)
			return nil, err
		}
		_dnsLog.Infof("mxs[%d] from dns for domain:%s", len(mxs), rDomain)
		if len(mxs) == 0 {
			mxs = []*net.MX{{Host: rDomain}}
		}
		du.Lock()
		du.MXs[rDomain] = mxs
		du.Unlock()
	}

	for _, mx := range mxs {
		c, err = tryConnect(lDomain, rDomain, mx)
		if err != nil {
			continue
		}
		return c, nil
	}

	du.Lock()
	delete(du.MXs, rDomain)
	du.Unlock()

	return nil, SMTPNoValidMX
}
