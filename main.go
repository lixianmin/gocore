/********************************************************************
created:    2018-10-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package main

import "github.com/lixianmin/gocore/service"

func main() {
	var server = service.NewTestServer()
	server.ListenAndServe()
}