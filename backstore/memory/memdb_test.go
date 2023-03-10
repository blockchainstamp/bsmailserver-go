package memory

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"github.com/emersion/go-msgauth/dmarc"
	"github.com/emersion/go-smtp"
	"net"
	"strings"
	"testing"
)

const mailHeaderString = "From: Joe SixPack <joe@football.example.com>\r\n" +
	"To: Suzie Q <suzie@shopping.example.net>\r\n" +
	"Subject: Is dinner ready?\r\n" +
	"Date: Fri, 11 Jul 2003 21:00:37 -0700 (PDT)\r\n" +
	"Message-ID: <20030712040037.46341.5F8J@football.example.com>\r\n"

const mailBodyString = "Hi.\r\n" +
	"\r\n" +
	"We lost the game. Are you hungry yet?\r\n" +
	"\r\n" +
	"Joe."
const dkim_domainkey = "v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAnEPsIS4dGLHyKp4zgaWQaGrBQ72598sWeND3wCgqkl8BC74XIzfpuXsZEmRB48acpH3thkoH4LTvInSmnL29vRjA970lXQWYcbPDoXZHY+YEojcy6cfJHKahJoXBleLsXBkoo72DYvDEAt5fKSY8Lrz5Nf8UkpTl4MGowwnyjS/LQ3jnGO5r1TfpfGfc+ZsXZVmehLI6c4ezjbcAw6QwlyTVMnMWf2ZGJ7IcyHJTpyY2O/QV/vzVrwa+47R/l+THrVupdb1CS2OQn2I5ld+9h/8Il7gBZfFXhstuFZd6AXlgG0w8U/5uHp1914ycNi8N5cFRUwLUXykFELsRc7aatf9ZQ+TTDL1GJ8y5rRiQlXcM2HJoBNrFJ4n88KGdTujfMSRerm1zhLTs+vmFiQ41LGEljsAJOM59lYs8LoBeQheWb/F3+EILp39Q9+YvqClWdhto6qLwOIWNptlKfRo+3txlTsn2knjQBZ9ih0+hmlhyuXvu0uWKkhU2XCQMf+wXAgMBAAE="
const dkim_priStr = "-----BEGIN PRIVATE KEY-----\n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\n/OF5FTPL/6bvHZpvWRwpqac=\n-----END PRIVATE KEY-----\n"
const mailString = mailHeaderString + "\r\n" + mailBodyString

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
	//tlsCfg, err := util.LoadServerTlsCnf("fullchain.pem", "privkey.pem")
	//if err != nil {
	//	t.Fatal(err)
	//}
	mxs, err := net.LookupMX("126.com")
	if err != nil {
		t.Fatal(err)
	}
	for _, mx := range mxs {
		fmt.Printf("\n%+v\n", mx)
		addr := fmt.Sprintf("%s:%d", mx.Host, util.DefaultSystemSmtpPort)
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

		//err = c.StartTLS(nil)
		//if err != nil {
		//	t.Fatal(err)
		//}
	}
}
