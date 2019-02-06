package gcp

import (
	"context"
	"net/http"
	"time"

	"github.com/puppetlabs/wash/datastore"
	"github.com/puppetlabs/wash/log"
	"github.com/puppetlabs/wash/plugin"
)

type project struct {
	name     string
	clientid string
	updated  time.Time
	clients  map[string]*service
}

// NewProject creates a new project with a collection of service clients.
func newProject(name string, clientid string, oauthClient *http.Client, cache *datastore.MemCache) (*project, error) {
	proj := &project{name, clientid, time.Now(), nil}
	services, err := newServices(name, proj.String(), oauthClient, cache)
	if err != nil {
		return nil, err
	}
	proj.clients = services
	return proj, nil
}

// Find service by name.
func (cli *project) Find(ctx context.Context, name string) (plugin.Node, error) {
	if svc, ok := cli.clients[name]; ok {
		log.Debugf("Found client %v", svc)
		return plugin.NewDir(svc), nil
	}
	return nil, plugin.ENOENT
}

// List all services as dirs.
func (cli *project) List(ctx context.Context) ([]plugin.Node, error) {
	log.Debugf("Listing %v clients in %v", len(cli.clients), cli)
	entries := make([]plugin.Node, 0, len(cli.clients))
	for _, svc := range cli.clients {
		entries = append(entries, plugin.NewDir(svc))
	}
	return entries, nil
}

// String returns a unique representation of the project.
func (cli *project) String() string {
	return cli.clientid + "/" + cli.Name()
}

// Name returns the project name.
func (cli *project) Name() string {
	return cli.name
}

// Attr returns attributes of the named service.
func (cli *project) Attr(ctx context.Context) (*plugin.Attributes, error) {
	return &plugin.Attributes{Mtime: cli.lastUpdate()}, nil
}

// Xattr returns a map of extended attributes.
func (cli *project) Xattr(ctx context.Context) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}

func (cli *project) closeServices(ctx context.Context) {
	log.Printf("Closing services for project %v", cli)
	for name, svc := range cli.clients {
		err := svc.close(ctx)
		if err != nil {
			log.Printf("Unable to close service %v in project %v: %v", name, cli.name, err)
		}
	}
}

func (cli *project) lastUpdate() time.Time {
	latest := cli.updated
	for _, svc := range cli.clients {
		if updated := svc.lastUpdate(); updated.After(latest) {
			latest = updated
		}
	}
	return latest
}
