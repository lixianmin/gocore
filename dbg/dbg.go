/********************************************************************
created:    2018-11-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package dbg

import (
	"net/http"
	"net/http/pprof"
)

type iServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

var _isDebugAllowed = func(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("Permission denied, please contact with lixianmin@live.cn"))
	return false
}

func EnableDebugger(mux iServeMux, isDebugAllowed func(w http.ResponseWriter, r *http.Request) bool) {
	if isDebugAllowed != nil {
		_isDebugAllowed = isDebugAllowed

		mux.HandleFunc("/debug/pprof/", index)
		mux.HandleFunc("/debug/pprof/cmdline", cmdline)
		mux.HandleFunc("/debug/pprof/profile", profile)
		mux.HandleFunc("/debug/pprof/symbol", symbol)
		mux.HandleFunc("/debug/pprof/trace", trace)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if _isDebugAllowed(w, r) {
		pprof.Index(w, r)
	}
}

func cmdline(w http.ResponseWriter, r *http.Request) {
	if _isDebugAllowed(w, r) {
		pprof.Cmdline(w, r)
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	if _isDebugAllowed(w, r) {
		pprof.Profile(w, r)
	}
}

func symbol(w http.ResponseWriter, r *http.Request) {
	if _isDebugAllowed(w, r) {
		pprof.Symbol(w, r)
	}
}

func trace(w http.ResponseWriter, r *http.Request) {
	if _isDebugAllowed(w, r) {
		pprof.Trace(w, r)
	}
}
