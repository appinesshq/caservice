// This program performs administrative tasks for the garage sale service.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/appinesshq/caservice/app/tooling/sales-admin/commands"
	"github.com/appinesshq/caservice/foundation/logger"
	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("ADMIN")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			log.Errorw("startup", "ERROR", err)
		}
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Args conf.Args
		Web  struct {
			APIHost       string `conf:"default:http://0.0.0.0:3000"`
			AuthTokenFile string `conf:"default:/tmp/.service/token"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Commands

	return processCommands(cfg.Args, cfg.Web.APIHost, cfg.Web.AuthTokenFile, log)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, host string, tokenfile string, log *zap.SugaredLogger) error {
	switch args.Num(0) {
	case "authenticate":
		if err := commands.Authenticate(host, args.Num(1), args.Num(2), tokenfile); err != nil {
			return fmt.Errorf("authenticate user: %w", err)
		}
	case "genkey":
		if err := commands.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}
	case "register":
		if err := commands.Register(host, args.Num(1), args.Num(2), args.Num(3)); err != nil {
			return fmt.Errorf("register user: %w", err)
		}
	default:
		fmt.Printf("authenticate: authenticate a user and store the token in %q\n", tokenfile)
		fmt.Println("genkey: generate a set of private/public key files")
		fmt.Println("register: register a new user")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
