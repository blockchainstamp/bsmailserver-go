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

func (du *DnsUtil) ValidSmtpCli(lDomain, rDomain string, tlsCfg *tls.Config) (*smtp.Client, error) {
	du.RLock()
	mxs, ok := du.MXs[rDomain]
	du.RUnlock()
	if !ok {
		_dnsLog.Info("no cached mx record", rDomain)
		mxs, err := net.LookupMX(rDomain)
		if err != nil {
			_dnsLog.Warn("LookupMX err:", err, rDomain)
			return nil, err
		}
		if len(mxs) == 0 {
			mxs = []*net.MX{{Host: rDomain}}
		}
		du.Lock()
		du.MXs[rDomain] = mxs
		du.Unlock()
	}

	for _, mx := range mxs {
		addr := fmt.Sprintf("%s:%d", mx.Host, DefaultSystemSmtpPort)
		c, err := smtp.Dial(addr)
		if err != nil {
			_dnsLog.Warn("dial err:", err, mx.Host)
			continue
		}
		err = c.Hello(lDomain)
		if err != nil {
			_dnsLog.Warn("say hello err:", err, lDomain, mx.Host)
			continue
		}
		if ok, _ := c.Extension("STARTTLS"); !ok {
			_dnsLog.Warn("server doesn't support STARTTLS", mx.Host)
			continue
		}
		err = c.StartTLS(tlsCfg)
		if err != nil {
			_dnsLog.Warn("start tls err:", err, mx.Host)
			continue
		}
		return c, nil
	}

	du.Lock()
	delete(du.MXs, rDomain)
	du.Unlock()

	return nil, SMTPNoValidMX
}
