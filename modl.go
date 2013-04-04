// -------
// modl.go ::: data models
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package main

import (
	//"labix.org/v2/mgo/bson"
	"net_kit/frms"
	"time"
)

type Coupon struct {
	Code, Discount string
	Active         bool
	User           string
	Timestamp      time.Time
}

type Admin struct {
	username, password string
}

func LoginForm() frms.Form {
	form := frms.Form{}
	form = frms.Form{
		Action: "/login",
		Button: "Login",
		Errors: make(map[string]string, 0),
		Inputs: []frms.Input{
			frms.Input{
				Name:     "username",
				Type:     "text",
				Holder:   "Username",
				Required: true,
			},
			frms.Input{
				Name:     "password",
				Type:     "password",
				Holder:   "Password",
				Required: true,
			},
		},
	}
	return form
}
