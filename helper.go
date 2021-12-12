package main

import (
	"fmt"
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/docker/docker-credential-helpers/credentials"
)

const (
	ghOAuthTokenConfigKey = "oauth_token"
	ghUserConfigKey       = "oauth_token"
	githubTokenUser       = "x-access-token"
)

var (
	ghFactory        = factory.New("")
	githubRegistries = []string{
		"docker.pkg.github.com",
		"ghcr.io",
	}
)

type Helper struct{}

var _ credentials.Helper = (*Helper)(nil)

func (h Helper) Add(_ *credentials.Credentials) error {
	return nil
}

func (h Helper) Delete(_ string) error {
	return nil
}

func (h Helper) Get(serverURL string) (string, string, error) {
	var ok bool
	for _, registry := range githubRegistries {
		if registry == serverURL {
			ok = true
			break
		}
	}
	if !ok {
		return "", "", fmt.Errorf("unknown registry: %s", serverURL)
	}

	config, err := ghFactory.Config()
	if err != nil {
		return "", "", err
	}

	host, err := config.DefaultHost()
	if err != nil {
		return "", "", err
	}

	token, source, err := config.GetWithSource(host, ghOAuthTokenConfigKey)
	if err != nil {
		return "", "", err
	}

	user := githubTokenUser
	if !strings.HasSuffix(source, "_TOKEN") {
		u, err := config.Get(host, ghUserConfigKey)
		if err != nil {
			return "", "", err
		}

		user = u
	}

	return user, token, nil
}

func (h Helper) List() (map[string]string, error) {
	config, err := ghFactory.Config()
	if err != nil {
		return nil, err
	}

	host, err := config.DefaultHost()
	if err != nil {
		return nil, err
	}

	_, source, err := config.GetWithSource(host, ghOAuthTokenConfigKey)
	if err != nil {
		return nil, err
	}

	user := githubTokenUser
	if !strings.HasSuffix(source, "_TOKEN") {
		u, err := config.Get(host, ghUserConfigKey)
		if err != nil {
			return nil, err
		}

		user = u
	}

	registries := make(map[string]string, len(githubRegistries))
	for _, registry := range githubRegistries {
		registries[registry] = user
	}

	return registries, nil
}
