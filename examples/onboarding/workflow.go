package onboarding

import (
	"log"

	"durable-execution-engine/engine"

	"golang.org/x/sync/errgroup"
)

func RunEmployeeOnboarding(ctx *engine.DurableContext) error {

	//Create emp record
	_, err := engine.AutoStep(ctx, func() (string, error) {
		log.Println("Creating employee record")
		return "employee-id-123", nil
	})
	if err != nil {
		return err
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := engine.AutoStep(ctx, func() (string, error) {
			log.Println("Provisioning laptop")
			return "laptop-ok", nil
		})
		return err
	})

	g.Go(func() error {
		_, err := engine.AutoStep(ctx, func() (string, error) {
			log.Println("Provisioning system access")
			return "access-ok", nil
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	_, err = engine.AutoStep(ctx, func() (string, error) {
		log.Println("Sending welcome email")
		return "email-sent", nil
	})

	return err
}
