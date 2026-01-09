// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/nextmn/cli-xdg"
	"github.com/nextmn/logrus-formatter/logger"

	"github.com/nextmn/upf/internal/app"
	"github.com/nextmn/upf/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func main() {
	logger.Init("NextMN-UPF")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	version := "Unknown version"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	app := &cli.Command{
		Name:                  "upf",
		Usage:                 "NextMN-UPF - Experimental 5G UPF",
		EnableShellCompletion: true,
		Authors: []any{
			"Louis Royer",
		},
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				TakesFile: true,
				Aliases:   []string{"c"},
				Usage:     "load configuration from `FILE`",
				// XXX: https://github.com/urfave/cli/issues/2244
				// Required:    true,
				DefaultText: "${XDG_CONFIG_DIRS}/nextmn-gnb-lite/config.yaml",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("CONFIG_FILE"),
					clixdg.ConfigFile("nextmn-gnb-lite/config.yaml"),
				),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// XXX: https://github.com/urfave/cli/issues/2244
			if cmd.String("config") == "" {
				logrus.Fatal("Required flag \"config\" not set")
			}

			conf, err := config.ParseConf(cmd.String("config"))
			if err != nil {
				logrus.WithContext(ctx).WithError(err).Fatal("Error loading config, exiting…")
			}
			if conf.Logger != nil {
				logrus.SetLevel(conf.Logger.Level)
			}

			if err := app.NewSetup(conf).Run(ctx); err != nil {
				logrus.WithError(err).Fatal("Error while running, exiting…")
			}
			return nil
		},
	}
	if err := app.Run(ctx, os.Args); err != nil {
		logrus.WithError(err).Fatal("Fatal error while running the application")
	}
}
