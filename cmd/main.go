package main

import (
	"flag"
	"log"

	"durable-execution-engine/engine"
	"durable-execution-engine/examples/onboarding"
)

var crashAfter = flag.Int("crash-after", -1, "Simulate crash after N steps")

func main() {
	flag.Parse()

	ctx, err := engine.NewDurableContext("employee-onboarding", "engine.db")
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.DB.Close()

	if *crashAfter > 0 {
		ctx.EnableCrashSimulation(*crashAfter)
	}

	if err := onboarding.RunEmployeeOnboarding(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Employee onboarding workflow completed")
}
