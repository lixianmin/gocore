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

var _isDebugAllowed = func(w http.ResponseWriter, r *http.Request) bool {
	_, _ = w.Write([]byte("Permission denied, please contact with lixianmin@live.cn"))
	return false
}

func EnableDebugger(rootDirectory string, handleFunc func(path string, handler func(http.ResponseWriter, *http.Request)), isDebugAllowed func(w http.ResponseWriter, r *http.Request) bool) {
	if nil == handleFunc {
		panic("handleFunc is nil")
	}

	if nil == isDebugAllowed {
		panic("isDebugAllowed is nil")
	}

	_isDebugAllowed = isDebugAllowed
	handleFunc(rootDirectory+"/debug/pprof/", index)
	handleFunc(rootDirectory+"/debug/pprof/cmdline", cmdline)
	handleFunc(rootDirectory+"/debug/pprof/profile", profile)
	handleFunc(rootDirectory+"/debug/pprof/symbol", symbol)
	handleFunc(rootDirectory+"/debug/pprof/trace", trace)
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
