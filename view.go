// -------
// view.go ::: view handlers
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package main

import (
	"labix.org/v2/mgo/bson"
	//"fmt"
	"net/http"
	"net_kit/data"
	"net_kit/frms"
	"net_kit/sess"
	"net_kit/tmpl"
	"net_kit/util"
	"strconv"
)

var (
	// init template store, session manager, and database wrapper
	ts = tmpl.NewTemplateStore("templates", "base.html")
	ss = sess.NewSessionStore("zoomcouponsid", sess.MIN*10)
	db = data.NewMgoWrapper("zoomenvelopes.com").SetDb("zoom")
)

func init() {
	// load templates and form structs into corrisponding caches
	ts.Load("login.html", "generator.html", "coupons.html", "err.html")
}

// get error
func errHandler(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(r.FormValue(":code"))
	ts.Render(w, "err.html", tmpl.M{"button": "generator", "code": code, "status": http.StatusText(code)})
}

// get login
func login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm()
	ts.Render(w, "login.html", tmpl.M{"button": "disabled", "form": form.Render(frms.INLINE)})
}

// post login
func loginAction(w http.ResponseWriter, r *http.Request) {
	form := LoginForm()
	if !form.IsValid(r) {
		ts.Render(w, "login.html", tmpl.M{"form": form.Render(frms.INLINE)})
		return
	}
	u, p := r.FormValue("username"), r.FormValue("password")
	if u == adm.username && p == adm.password {
		s := ss.NewSession(w, r)
		s.Set("authd", []string{"y"})
		http.Redirect(w, r, "/generator", 301)
		return
	}
	form.Errors["username"] = "invalid username or password"
	ts.Render(w, "login.html", tmpl.M{"button": "disabled", "form": form.Render(frms.INLINE)})
	return
}

// get logout
func logout(w http.ResponseWriter, r *http.Request) {
	ss.DelSession(w, r)
	http.Redirect(w, r, "/", 301)
	return
}

// get generator
func generator(w http.ResponseWriter, r *http.Request) {
	s := ss.GetSession(w, r)
	ok := s.Get("authd")[0]
	if ok != "y" {
		ss.DelSession(w, r)
		http.Redirect(w, r, "/", 301)
		return
	}
	ts.Render(w, "generator.html", tmpl.M{"button": "logout"})
	return
}

// post generator
func generatorAction(w http.ResponseWriter, r *http.Request) {
	s := ss.GetSession(w, r)
	ok := s.Get("authd")[0]
	if ok != "y" {
		ss.DelSession(w, r)
		http.Redirect(w, r, "/", 301)
		return
	}
	d, q := r.FormValue("discount"), r.FormValue("quantity")
	if d != "" && q != "" {
		var coupons []interface{}
		n, _ := strconv.Atoi(q)
		for i := 0; i < n; i++ {
			coupons = append(coupons, Coupon{Code: util.Md5()[:10], Discount: d, Active: true})
		}
		db.SetC("coupon").Insert(coupons...)
		http.Redirect(w, r, "/coupons", 301)
		return
	}
	http.Redirect(w, r, "/error/405", 301)
	return
}

// get coupons
func coupons(w http.ResponseWriter, r *http.Request) {
	s := ss.GetSession(w, r)
	ok := s.Get("authd")[0]
	if ok != "y" {
		ss.DelSession(w, r)
		http.Redirect(w, r, "/", 301)
		return
	}
	var coupons []Coupon
	db.SetC("coupon").Return(0, bson.M{"active": true}, &coupons)
	ts.Render(w, "coupons.html", tmpl.M{"button": "generator", "coupons": coupons})
	return
}
