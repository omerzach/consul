package xds

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/hashicorp/consul/agent/structs"
)

func firstHealthyTarget(
	targetConfigs map[structs.DiscoveryTarget]structs.DiscoveryTargetConfig,
	targetHealth map[structs.DiscoveryTarget]structs.CheckServiceNodes,
	primary structs.DiscoveryTarget,
	secondary []structs.DiscoveryTarget,
) structs.DiscoveryTarget {
	series := make([]structs.DiscoveryTarget, 0, len(secondary)+1)
	series = append(series, primary)
	series = append(series, secondary...)

	for _, target := range series {
		targetConfig := targetConfigs[target]

		endpoints, ok := targetHealth[target]
		if !ok {
			continue
		}
		for _, ep := range endpoints {
			healthStatus, _ := calculateEndpointHealthAndWeight(ep, targetConfig.Subset.OnlyPassing)
			if healthStatus == envoycore.HealthStatus_HEALTHY {
				return target
			}
		}
	}

	return primary // if everything is broken just use the primary for now
}
