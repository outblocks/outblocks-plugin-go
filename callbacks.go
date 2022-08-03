package plugin

import apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"

func DefaultRegistryApplyCallback(stream apiv1.DeployPluginService_ApplyServer) func(*apiv1.ApplyAction) {
	return func(a *apiv1.ApplyAction) {
		_ = stream.Send(&apiv1.ApplyResponse{
			Response: &apiv1.ApplyResponse_Action{
				Action: &apiv1.ApplyActionResponse{
					Actions: []*apiv1.ApplyAction{a},
				},
			},
		})
	}
}

func DefaultRegistryApplyDNSCallback(stream apiv1.DNSPluginService_ApplyDNSServer) func(*apiv1.ApplyAction) {
	return func(a *apiv1.ApplyAction) {
		_ = stream.Send(&apiv1.ApplyDNSResponse{
			Response: &apiv1.ApplyDNSResponse_Action{
				Action: &apiv1.ApplyActionResponse{
					Actions: []*apiv1.ApplyAction{a},
				},
			},
		})
	}
}

func DefaultRegistryApplyMonitoringCallback(stream apiv1.MonitoringPluginService_ApplyMonitoringServer) func(*apiv1.ApplyAction) {
	return func(a *apiv1.ApplyAction) {
		_ = stream.Send(&apiv1.ApplyMonitoringResponse{
			Response: &apiv1.ApplyMonitoringResponse_Action{
				Action: &apiv1.ApplyActionResponse{
					Actions: []*apiv1.ApplyAction{a},
				},
			},
		})
	}
}
