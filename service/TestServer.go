/********************************************************************
created:    2018-11-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package service

import (
	"encoding/json"
	"github.com/lixianmin/gocore/dbg"
	"net/http"
)

type TestServer struct {
}

func NewTestServer() *TestServer {
	var server = &TestServer{}
	return server
}

func (server *TestServer) ListenAndServe() {
	//  处理debug消息
	var mux = http.NewServeMux()
	dbg.EnableDebugger(mux.HandleFunc, func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("h") == "1" {
			return true
		}

		//w.Write([]byte("Permission denied, please contact with lixianmin@live.cn"))
		return true
	})

	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]string{"status": "UP", "notes": "some description"}
		response, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	var err = http.ListenAndServe(":8444", mux)
	if err != nil {
		//console.Error("[listenAndServe()] failed with err= %q", err)
	}
}
