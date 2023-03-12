module github.com/blockchainstamp/bsmailserver-go

go 1.19

require (
	github.com/blockchainstamp/go-mail-proxy v1.2.2
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-message v0.15.0
	github.com/emersion/go-msgauth v0.6.6
	github.com/emersion/go-smtp v0.16.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/emersion/go-sasl v0.0.0-20220912192320-0145f2c60ead // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

replace github.com/emersion/go-smtp => github.com/blockchainstamp/go-smtp v0.0.0-20230108191019-90d596c5fb00
