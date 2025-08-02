package gottylib

import (
	"context"
	"fmt"
	"github.com/kernel-punk/gotty/backend/localcommand"
	"github.com/kernel-punk/gotty/server"
	"github.com/kernel-punk/gotty/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Run(runParameters RunParameters) error {

	//runParameters := RunParameters{
	//
	//	cmd:  "htop",
	//	args: []string{},
	//	ssl:  false,
	//}

	hostname, _ := os.Hostname()

	appOptions := &server.Options{}
	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return err
	}

	backendOptions := &localcommand.Options{}
	if err := utils.ApplyDefaultValues(backendOptions); err != nil {
		return err
	}

	factory, err := localcommand.NewFactory(runParameters.Cmd, runParameters.Args, backendOptions)
	if err != nil {
		return err
	}

	appOptions.TitleVariables = map[string]interface{}{
		"command":  runParameters.Cmd,
		"argv":     runParameters.Args,
		"hostname": hostname,
	}

	appOptions.PermitWrite = true

	appOptions.Port = strconv.Itoa(runParameters.Port)

	srv, err := server.New(factory, appOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	gCtx, gCancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)
	go func() {
		errs <- srv.Run(ctx, server.WithGracefullContext(gCtx))
	}()
	err = waitSignals(errs, cancel, gCancel)

	if err != nil && err != context.Canceled {
		return err
	}

	return nil
}

func waitSignals(errs chan error, cancel context.CancelFunc, gracefullCancel context.CancelFunc) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-errs:
		return err

	case s := <-sigChan:
		switch s {
		case syscall.SIGINT:
			gracefullCancel()
			fmt.Println("C-C to force close")
			select {
			case err := <-errs:
				return err
			case <-sigChan:
				fmt.Println("Force closing...")
				cancel()
				return <-errs
			}
		default:
			cancel()
			return <-errs
		}
	}
}
