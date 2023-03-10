package main

import (
	"fmt"
	bp "github.com/blockchainstamp/bsmailserver-go"
	"github.com/blockchainstamp/bsmailserver-go/cfg"
	"github.com/blockchainstamp/bsmailserver-go/service"
	"github.com/blockchainstamp/go-mail-proxy/utils/fdlimit"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "bsmail",
	Short: "mail service support blockchain stamp",
	Long:  `TODO::.`,
	Run:   mainRun,
}

var (
	version   bool
	withStamp bool
	mailBase  string
)

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().BoolVarP(&withStamp, "stamp", "s", false,
		"run server with stamp logic")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false,
		"version of current system")
	rootCmd.Flags().StringVarP(&mailBase, "base", "b", "",
		"base directory of current system")
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func initSystem() error {

	if err := os.Setenv("GODEBUG", "netdns=go"); err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	limit, err := fdlimit.Maximum()
	if err != nil {
		return fmt.Errorf("failed to retrieve file descriptor allowance:%s", err)
	}
	_, err = fdlimit.Raise(uint64(limit))
	if err != nil {
		return fmt.Errorf("failed to raise file descriptor allowance:%s", err)
	}
	return nil
}

func mainRun(_ *cobra.Command, _ []string) {
	if version {
		logVersion()
		return
	}
	var (
		err error
	)

	err = initSystem()
	if err != nil {
		fmt.Println("init fd limit of system failed:=>", err)
		return
	}

	if err = cfg.PrepareConfig(mailBase); err != nil {
		fmt.Println("load system config failed:=>", err)
		return
	}
	cfg.SetStampTag(withStamp)
	service.Inst().Run()
	waitSignal()
	service.Inst().ShutDown()
}

func logVersion() {
	fmt.Println("\n==================================================")
	fmt.Printf("Version:\t %s\n", bp.Version)
	fmt.Printf("Build:\t%s\n", bp.BuildTime)
	fmt.Printf("Commit:\t%s\n", bp.Commit)
	fmt.Println("==================================================")
}

func waitSignal() {
	sigCh := make(chan os.Signal, 1)
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>bsmail start at pid(%s)<<<<<<<<<<\n", pid)

	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		os.Kill,
	)

	//TODO:: process system signal
	for sig := range sigCh {
		fmt.Printf("\n>>>>>>>>>>bsmail[%s] finished(%s)<<<<<<<<<<\n", pid, sig)
		switch sig {
		case syscall.SIGHUP:
			fmt.Println(">>>>>>>>>>hangup<<<<<<<<<<")
			continue
		default:
			return
		}
	}
}
