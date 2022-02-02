package opsgenie

import (
	"context"
	"fmt"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type HeartbeatGenerator struct {
	OpsgenieService
}

func (g *HeartbeatGenerator) InitResources() error {
	client, err := g.HeartbeatClient()
	if err != nil {
		return err
	}

	var heartbeats []heartbeat.Heartbeat

	result, err := func() (*heartbeat.ListResult, error) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFunc()

		return client.List(ctx)
	}()

	if err != nil {
		return err
	}

	heartbeats = append(heartbeats, result.Heartbeats...)

	g.Resources = g.createResources(heartbeats)
	return nil
}

func (g *HeartbeatGenerator) createResources(heartbeats []heartbeat.Heartbeat) []terraformutils.Resource {
	var resources []terraformutils.Resource

	for _, h := range heartbeats {
		resources = append(resources, terraformutils.NewResource(
			fmt.Sprintf("Heartbeat-%s", h.Name),
			h.Name,
			"opsgenie_heartbeat",
			g.ProviderName,
			map[string]string{},
			[]string{},
			map[string]interface{}{},
		))
	}

	return resources
}
