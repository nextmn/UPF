// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"

	"github.com/nextmn/upf/internal/constants"

	pfcp_networking "github.com/nextmn/go-pfcp-networking/pfcp"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func (s *Setup) createTUNInterface(ctx context.Context) error {
	if s.tunInterface == nil {
		return fmt.Errorf("TUN interface has not been created")
	}
	var gtpEntity netip.Addr
	if len(s.config.Gtpu.GTPUProtocolEntities) > 0 {
		gtpEntity = s.config.Gtpu.GTPUProtocolEntities[0].Addr
	}
	go func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				packet := make([]byte, constants.MTU_GTP_TUN)
				n, err := s.tunInterface.Read(packet)
				if err != nil {
					return err
				}
				go func(ctx context.Context, packet []byte, db *FARAssociationDB, tuniface *water.Interface, pfcpServer *pfcp_networking.PFCPEntityUP) {
					err := ipPacketHandler(ctx, gtpEntity, packet, db, tuniface, pfcpServer)
					if err != nil {
						logrus.WithError(err).Debug("Drop packet")
					}
				}(ctx, packet[:n], s.farUconnDb, s.tunInterface, s.pfcpServer)
			}
		}
	}(ctx)
	return nil
}

func (s *Setup) createDLRoutes(ctx context.Context) error {
	if s.tunInterface == nil {
		return fmt.Errorf("TUN interface has not been created")
	}
	for _, ue := range s.config.DNNList {
		err := runIP(ctx, "route", "add", ue.Cidr, "dev", s.tunInterface.Name(), "proto", "static")
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"prefix": ue.Cidr}).Error("Cannot create Uplink route for this prefix")
			return err
		}
	}
	return nil
}

func (s *Setup) removeTun(ctx context.Context) error {
	if s.tunInterface == nil {
		return nil
	}
	err := runIP(ctx, "link", "del", s.tunInterface.Name())
	if nil != err {
		logrus.WithError(err).WithFields(logrus.Fields{"interface": s.tunInterface.Name()}).Error("Unable to delete interface")
		return err
	}
	return nil
}
func (s *Setup) createTun(ctx context.Context) error {
	config := water.Config{
		DeviceType: water.TUN,
	}
	if s.config.Gtpu.TunInterfaceName != nil {
		config.Name = *s.config.Gtpu.TunInterfaceName
	}
	iface, err := water.New(config)
	if nil != err {
		logrus.WithError(err).Error("Unable to allocate TUN interface")
		return err
	}
	err = runIP(ctx, "link", "set", "dev", iface.Name(), "mtu", strconv.Itoa(constants.MTU_GTP_TUN))
	if nil != err {
		logrus.WithError(err).WithFields(logrus.Fields{
			"mtu":       constants.MTU_GTP_TUN,
			"interface": iface.Name(),
		}).Error("Unable to set MTU")
		return err
	}
	err = runIP(ctx, "link", "set", "dev", iface.Name(), "up")
	if nil != err {
		logrus.WithError(err).WithFields(logrus.Fields{
			"interface": iface.Name(),
		}).Error("Unable to set interface up")
		return err
	}
	s.tunInterface = iface

	err = runIPTables(ctx, "-A", "OUTPUT", "-o", iface.Name(), "-p", "icmp", "--icmp-type", "redirect", "-j", "DROP")
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"interface": iface.Name()}).Error("Error while setting iptable rule to drop icmp redirects")
		return err
	}
	err = runIP6Tables(ctx, "-A", "OUTPUT", "-o", iface.Name(), "-p", "icmpv6", "--icmpv6-type", "redirect", "-j", "DROP")
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"interface": iface.Name()}).Error("Error while setting iptable rule to drop icmpv6 redirects")
		return err
	}
	return nil
}
