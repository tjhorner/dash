package main

import "github.com/tjhorner/dash/util"

var envListenAddr = util.GetEnv("DASH_LISTEN_ADDRESS", ":3000")
