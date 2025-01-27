// A generated module for Rewild functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/rewild/internal/dagger"
	"fmt"
)

type Rewild struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Rewild) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Rewild) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

const registryHost = "sjc.vultrcr.com"
const registry = "https://" + registryHost + "/oosa"

func (m *Rewild) BuildProduct(ctx context.Context,
	source *dagger.Directory, image string, dockerfile string,
	tag string, user string, password *dagger.Secret) ([]string, error) {

	var platforms = []dagger.Platform{
		"linux/amd64", // a.k.a. x86_64
		"linux/arm64", // a.k.a. aarch64
	}

	var buildTags []string
	if tag == "" {
		buildTags = []string{"latest"}
	} else {
		buildTags = []string{tag, "latest"}
	}

	var container *dagger.Container
	refs := make([]string, len(buildTags))
	for i, buildTag := range buildTags {
		imageRepo := fmt.Sprintf("%s/oosa/%s:%s", registryHost, image, buildTag)
		platformVariants := make([]*dagger.Container, 0, len(platforms))
		for _, platform := range platforms {
			container = source.
				DockerBuild(dagger.DirectoryDockerBuildOpts{
					Platform:   platform,
					Dockerfile: dockerfile,
				})
			platformVariants = append(platformVariants, container)

		}

		ref, err := dag.Container().
			WithRegistryAuth(registry, user, password).
			Publish(ctx, imageRepo, dagger.ContainerPublishOpts{
				PlatformVariants: platformVariants,
			})
		if err != nil {
			return nil, err
		}
		refs[i] = ref
	}

	return refs, nil
}
