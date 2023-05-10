package main

import (
	"fmt"
	"log"
	"obAggregator/actors/delta"
	"obAggregator/actors/superuser"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	startServices()
}

func startServices() {
	gracefulStop := make(chan os.Signal, 8)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV, os.Interrupt)
	go startSupervisors(gracefulStop)
	<-gracefulStop
	os.Exit(1)
}

func startSupervisors(gracefulStop chan os.Signal) {
	var actorPID *actor.PID
	system := actor.NewActorSystem()
	serviceToStart := os.Args[len(os.Args)-1]
	port := 9082
	if serviceToStart == "--delta" {
		port = 9081
	}
	r := remote.NewRemote(system, remote.Configure("0.0.0.0", port,
		remote.WithAdvertisedHost(fmt.Sprintf("0.0.0.0:%v", port))))
	r.Start()
	decider := func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
	supervisor := actor.NewOneForOneStrategy(10, 1000, decider)
	rootContext := system.Root
	if serviceToStart == "--delta" {
		props := actor.
			PropsFromProducer(delta.NewDeltaWorker("BTC_USDT,ETH_USDT", "BINANCE,DERIBIT", supervisor), actor.WithSupervisor(supervisor))
		pid, err := rootContext.SpawnNamed(props, "DELTA")
		if err != nil {
			log.Println("failed to start child-actor for bot", err)
		}
		actorPID = pid
	} else if serviceToStart == "--enduser" {
		props := actor.
			PropsFromProducer(superuser.NewDSuperUser(supervisor), actor.WithSupervisor(supervisor))
		pid, err := rootContext.SpawnNamed(props, "SUPERUSER")
		if err != nil {
			log.Println("failed to start child-actor for bot", err)
		}
		actorPID = pid
	}
	<-gracefulStop
	rootContext.Poison(actorPID)
	time.Sleep(1 * time.Second)
}
