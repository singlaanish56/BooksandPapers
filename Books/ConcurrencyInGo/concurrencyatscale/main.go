package main

import (
	"github.com/singlaanish56/Books/ConcurrencyInGo/concurrencyatscale/ratelimiting"
)

func main(){
	// errorpropogation.RunErrorProp()
	// errorpropogation.RunErrorPropBetter()
	//heartbeat.BasicHeartbeart()
	//heartbeat.HeartbeartThatBreaks()
	//heartbeat.UnitWorkHeartbeat()
	//replicatedRequests.RunDoWork()
	//ratelimiting.ClientSideNoRateLimiter()
	//ratelimiting.ClientSideSingleRateLimiter()
	ratelimiting.ClientSideMultiLimiter()
}