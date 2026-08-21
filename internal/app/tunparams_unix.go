// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

//go:build linux || darwin

package app

import "github.com/songgao/water"

func platformSpecificParams(name string) water.PlatformSpecificParams {
	return water.PlatformSpecificParams{
		Name: name,
	}
}
