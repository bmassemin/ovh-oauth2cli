package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cli/browser"
	"github.com/int128/oauth2cli"
	"github.com/int128/oauth2cli/oauth2params"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

//go:embed cert/server.crt
var cert []byte

//go:embed cert/server.key
var key []byte

type cmdOptions struct {
	clientID        string
	clientSecret    string
	scopes          string
	region          string
	localServerCert string
	localServerKey  string
	port            int
}

func main() {
	log.SetOutput(os.Stderr)

	var o cmdOptions
	flag.StringVar(&o.clientID, "client-id", "", "OAuth Client ID")
	flag.StringVar(&o.clientSecret, "client-secret", "", "OAuth Client Secret")
	flag.StringVar(&o.scopes, "scopes", "all", "Scopes to request, comma separated")
	flag.StringVar(&o.region, "region", "eu", "Region to use")
	flag.StringVar(&o.localServerCert, "server-crt", "", "Path to a certificate file for the local server")
	flag.StringVar(&o.localServerKey, "server-key", "", "Path to a key file for the local server")
	flag.IntVar(&o.port, "port", 8080, "Port for the local server")
	flag.Parse()

	if o.clientID == "" || o.clientSecret == "" {
		log.Print(`You need to set oauth2 credentials.
Create a service account at https://eu.api.ovh.com/console/?section=%2Fme&branch=v1#post-/me/api/oauth2/client with AUTHORIZATION_CODE flow
Then set the following options:`)
		flag.PrintDefaults()
		os.Exit(1)
		return
	}

	if o.localServerCert == "" || o.localServerKey == "" {
		cleanup := setupEmbeddedCerts(&o)
		defer cleanup()
	}

	regionEndpoints := map[string]string{
		"eu": "https://www.ovh.com",
		"ca": "https://ca.ovh.com",
	}
	endpoint, ok := regionEndpoints[o.region]
	if !ok {
		log.Fatalf("region %s is not supported", o.region)
	}

	pkce, err := oauth2params.NewPKCE()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	ready := make(chan string, 1)
	defer close(ready)
	cfg := oauth2cli.Config{
		OAuth2Config: oauth2.Config{
			ClientID:     o.clientID,
			ClientSecret: o.clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s/auth/oauth2/authorize", endpoint),
				TokenURL: fmt.Sprintf("%s/auth/oauth2/token", endpoint),
			},
			Scopes: []string{"all"},
		},
		AuthCodeOptions:        pkce.AuthCodeOptions(),
		TokenRequestOptions:    pkce.TokenRequestOptions(),
		LocalServerReadyChan:   ready,
		LocalServerCertFile:    o.localServerCert,
		LocalServerKeyFile:     o.localServerKey,
		LocalServerBindAddress: []string{fmt.Sprintf("localhost:%d", o.port)},
		Logf:                   log.Printf,
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			if err := browser.OpenURL(url); err != nil {
				log.Printf("could not open the browser: %s", err)
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		token, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}
		setEnvVar(o.region, token.AccessToken)
		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalf("authorization error: %s", err)
	}
}

func setEnvVar(region string, accessToken string) {
	confFile := fmt.Sprintf(`[default]
endpoint=ovh-%s

[ovh-%s]
access_token=%s`, region, region, accessToken)
	os.WriteFile("ovh.conf", []byte(confFile), 0644)
}

func setupEmbeddedCerts(o *cmdOptions) (cleanup func()) {
	tmpDir := os.TempDir()
	o.localServerCert = filepath.Join(tmpDir, "server.crt")
	o.localServerKey = filepath.Join(tmpDir, "server.key")
	if err := os.WriteFile(o.localServerCert, cert, 0644); err != nil {
		log.Fatalf("could not write certificate file: %s", err)
	}
	if err := os.WriteFile(o.localServerKey, key, 0644); err != nil {
		log.Fatalf("could not write key file: %s", err)
	}
	return func() {
		if err := os.Remove(o.localServerCert); err != nil {
			log.Printf("could not remove certificate file: %s", err)
		}
		if err := os.Remove(o.localServerKey); err != nil {
			log.Printf("could not remove key file: %s", err)
		}
	}
}
