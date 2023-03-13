package util

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/emersion/go-msgauth/dmarc"
	"github.com/emersion/go-smtp"
	"gopkg.in/gomail.v2"
	"net"
	"strings"
	"testing"
	"time"
)

var (
	host = ""
)

func init() {
	flag.StringVar(&host, "host", "smtp.126.com", "")
}

func TestTxtLookup(t *testing.T) {
	txts, err := net.LookupTXT("_dmarc." + "simplenets.org")
	if err != nil {
		t.Fatal(err)
	}
	txt := strings.Join(txts, "")
	rec, err := dmarc.Parse(txt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", rec)
	txts, err = net.LookupTXT("dkim" + "._domainkey." + "simplenets.org")
	txt = strings.Join(txts, "")
	fmt.Println(txt)
}

func TestMXInfo(t *testing.T) {

	mxs, err := net.LookupMX("simplenets.org.")
	if err != nil {
		t.Fatal(err)
	}
	for _, mx := range mxs {
		fmt.Printf("\n%+v\n", mx)
		addr := fmt.Sprintf("%s:%d", mx.Host, DefaultSystemSmtpPort)
		c, err := smtp.Dial(addr)
		if err != nil {
			t.Fatal(err)
		}
		err = c.Hello("simplenets.org")
		if err != nil {
			t.Fatal(err)
		}
		if ok, _ := c.Extension("STARTTLS"); !ok {
			t.Fatal(err)
		}
	}
}

func TestSimpleDial(t *testing.T) {
	conn, err := net.DialTimeout("tcp", host+":25", time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	cli, err := smtp.NewClient(conn, host)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.Hello("localhost")
	if err != nil {
		t.Fatal(err)
	}
	if ok, _ := cli.Extension("STARTTLS"); !ok {
		t.Fatal("must support STARTTLS")
	}
	tlsCfg := &tls.Config{ServerName: host}
	if err := cli.StartTLS(tlsCfg); err != nil {
		t.Fatal(err)
	}
}

func TestTlsDial(t *testing.T) {
	conn, err := net.DialTimeout("tcp", host+":465", time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	tlsCfg := &tls.Config{ServerName: host}
	conn = tls.Client(conn, tlsCfg)
	cli, err := smtp.NewClient(conn, host)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.Hello("localhost")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoMail(t *testing.T) {
	gomail.NewDialer("smtp.simplenets.org", 25, "", "")
}

func TestDnsQuery(t *testing.T) {
	mxs, err := net.LookupMX(host)
	if err != nil {
		t.Fatal(err)
	}
	for _, mx := range mxs {
		c, err := tryConnect("localhost", host, mx)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(c.TLSConnectionState())
	}
}

var (
	now = time.Now().Format("2006-01-02 15:04:05")

	tEmail  = "midoks@163.com"
	fEmail  = "admin@cachecha.com"
	content = fmt.Sprintf("From: <%s>\r\nSubject: Hello imail[%s]\r\nTo: <%s>\r\n\r\nHi! yes is test. imail ok?", fEmail, now, tEmail)
)

func TestTryConn(t *testing.T) {
	var addr = host + ":25"
	conn, err := net.DialTimeout("tcp4", addr, MailMTATimeOut)
	if err != nil {
		t.Fatal(err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Hello("localhost")
	if err != nil {
		t.Fatal(err)
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: host}
		err = c.StartTLS(tlsCfg)
		if err != nil {
			t.Fatal(err)
		}
	}
	if ok, _ := c.Extension("AUTH"); !ok {
		t.Fatal("smtp: server doesn't support AUTH")
	}

	//if err = c.Mail(from); err != nil {
	//	return err
	//}
	//for _, addr := range to {
	//	if err = c.Rcpt(addr); err != nil {
	//		return err
	//	}
	//}
	//w, err := c.Data()
	//
	//if err != nil {
	//	return err
	//}
	//_, err = w.Write(msg)
	//
	//if err != nil {
	//	return err
	//}
	//err = w.Close()
	//if err != nil {
	//	return err
	//}
	//return c.Quit()
}
