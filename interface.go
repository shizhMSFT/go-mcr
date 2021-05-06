package mcr

import "context"

type Client interface {
	Repositories(ctx context.Context) ([]string, error)
	Tags(ctx context.Context, repo string) ([]string, error)
	Manifest(ctx context.Context, repo, ref string) (string, []byte, error)
}
