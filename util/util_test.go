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
