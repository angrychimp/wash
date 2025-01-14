package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/puppetlabs/wash/activity"
	"github.com/puppetlabs/wash/plugin"
)

type containersDir struct {
	plugin.EntryBase
	client *client.Client
}

func newContainersDir(client *client.Client) *containersDir {
	containersDir := &containersDir{
		EntryBase: plugin.NewEntry("containers"),
	}
	containersDir.client = client
	return containersDir
}

func (cs *containersDir) Schema() *plugin.EntrySchema {
	return plugin.NewEntrySchema(cs, "containers").IsSingleton()
}

func (cs *containersDir) ChildSchemas() []*plugin.EntrySchema {
	return []*plugin.EntrySchema{
		(&container{}).Schema(),
	}
}

// List
func (cs *containersDir) List(ctx context.Context) ([]plugin.Entry, error) {
	containers, err := cs.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	activity.Record(ctx, "Listing %v containers in %v", len(containers), cs)
	keys := make([]plugin.Entry, len(containers))
	for i, inst := range containers {
		keys[i] = newContainer(inst, cs.client)
	}
	return keys, nil
}
