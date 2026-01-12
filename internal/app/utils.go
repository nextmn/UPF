// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func runIP(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ip", args...)
	cmd.Env = []string{}
	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"command":   cmd.Args[0],
			"arguments": args,
		}).WithError(err).Error("Error while running command")
		return err
	}
	return nil
}

func runIPTables(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "iptables", args...)
	cmd.Env = []string{}
	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"command":   cmd.Args[0],
			"arguments": args,
		}).WithError(err).Error("Error while running command")
		return err
	}
	return nil
}

func runIP6Tables(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ip6tables", args...)
	cmd.Env = []string{}
	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"command":   cmd.Args[0],
			"arguments": args,
		}).WithError(err).Error("Error while running command")
		return err
	}
	return nil
}
