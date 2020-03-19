package ignition

import (
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
)

type FakeBackend struct {
}

func (factory *FakeBackend) getIngitionConfigTemplate(node *Node) (*ignTypes.Config, error) {
	out := &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: IngitionSchemaVersion,
		},
	}
	return out, nil
}

func (factory *FakeBackend) applyConfig(config *ignTypes.Config) (*ignTypes.Config, error) {
	return config, nil
}
