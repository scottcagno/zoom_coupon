// -------
// cont.go ::: http controller
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package main

import (
	"net_kit/web"
)

var (
	mux = web.NewMultiplexer()
	srv = web.NewWebServer()
	adm = Admin{"admin", "env52nei"}
)

func main() {

	// forward to login
	mux.Forward("/", "/login")

	// logging in and out
	mux.Get("/login", login)
	mux.Post("/login", loginAction)
	mux.Get("/logout", logout)

	// coupon generator
	mux.Get("/generator", generator)
	mux.Post("/generator", generatorAction)

	// overview
	mux.Get("/coupons", coupons)

	// errors
	mux.Get("/error/:code", errHandler)

	// static
	mux.Static("/static/", "static")

	// serve
	srv.Serve(":8080", mux)
}
