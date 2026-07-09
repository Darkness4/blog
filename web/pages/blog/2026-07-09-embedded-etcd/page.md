---
title: Using embedded etcd as distributed local store.
description: Easy high availability for stateful services.
tags: ["go", "distributed", "programming", "etcd"]
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

## Introduction

There isn't many documentation about embedding etcd as a distributed local
store, so I wanted to write a simple article about it.

There are some situations where you do not want to add more complexity to a
deployment by deploying an additional remote store. Sometimes, you feel like the
level of availability of your application doesn't matter and should match the
availability of the store. Or, simply, you have an actual stateful application,
but want to distribute it.

In previous articles, I talked about how etcd can be used to delegate the state
of the application to a distributed store to achieve high availability. Today,
I will present a mix of both world: a distributed stateful application with
delegated state.

## Quick remainder of etcd

Etcd is a highly consistent and reliable distributed key-value store. It uses
the Raft consensus algorithm to ensure that the state of the store is always
consistent across the cluster.

To create a cluster, replicas of etcd are deployed, and a persistent volume is attached to each
replica. Each replica knows the addresses of the other replicas.

It is a very popular distributed store for stateful applications like the
Kubernetes control plane.

## Basic etcd cluster settings

First, let's get familiar on how's etcd is deployed. You need these parameters:

- **Storage settings**:
  - `data-dir`: The directory where the data is stored.
- **Network settings**:
  - `listen-client-urls`: The addresses used to listen for client requests. Since
    we are using the embedded etcd, we will use `http://localhost:2379`. We don't
    want to expose to external traffic.
  - `listen-peer-urls`: The addresses used to listen for peer requests. This one
    needs to be public for the peers to connect to each other. By default, we set
    `http://0.0.0.0:2380`.
- **Cluster settings**:
  - `name`: The identifier of the etcd instance, used for clustering.
  - `advertise-client-urls`: The addresses used to advertise to the rest of the
    cluster. By default, we set `http://localhost:2379`.
  - `initial-advertise-peer-urls`: The addresses used to advertise the peer
    address to the rest of the cluster. This is more delicate, you should use the
    public address of the peer. By default, we set to `http://localhost:2380` (for
    single node cluster).
  - `initial-cluster`: The initial cluster configuration. This is a comma separated
    list of peer addresses. For example, `infra0=http://localhost:2380`.
  - `initial-cluster-state`: The initial state of the cluster. By default, we set
    to `new`. If you need to add new members, you set the value to `existing`.
  - `initial-cluster-token`: The initial token of the cluster. It should be unique
    to avoid conflict.

## Using embedded etcd as a distributed local store

### Design of the application

For the sake of the examples, we'll create a service that achieve this simple
use case:

"A highly available remote configuration store for users".

Users will be able to store JSON data in the store, and retrieve it.

### Bootstrapping the application

I'm going to use `urfav` to manage the CLI. Let's handle the flags specified
above:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

var (
	version = "dev"

	etcdName                     string
	etcdListenClientURLs         []string
	etcdListenPeerURLs           []string
	etcdAdvertiseClientURLs      []string
	etcdInitialAdvertisePeerURLs []string
	etcdDataDir                  string
	etcdInitialClusterState      string
	etcdInitialCluster           string
	etcdInitialClusterToken      string
)

var app = &cli.Command{
	Name:        "Distributed user store.",
	Description: `Remote HA user store.`,
	Suggest:     true,
	Version:     version,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "etcd.name",
			Usage:       "Name of etcd instance",
			Value:       "etcd-node",
			Destination: &etcdName,
			Sources: cli.EnvVars(
				"ETCD_NAME",
			),
		},
		&cli.StringSliceFlag{
			Name:        "etcd.listen-client-urls",
			Usage:       "Listen address for client traffic. You should set to localhost or 127.0.0.1.",
			Value:       []string{"http://localhost:2379"},
			Destination: &etcdListenClientURLs,
			Sources: cli.EnvVars(
				"ETCD_LISTEN_CLIENT_URLS",
			),
		},
		&cli.StringSliceFlag{
			Name:        "etcd.listen-peer-urls",
			Usage:       "Listen address for peer traffic.",
			Value:       []string{"http://0.0.0.0:2380"},
			Destination: &etcdListenPeerURLs,
			Sources: cli.EnvVars(
				"ETCD_LISTEN_PEER_URLS",
			),
		},
		&cli.StringSliceFlag{
			Name:        "etcd.advertise-client-urls",
			Usage:       "Advertise address for client traffic. You should set to localhost.",
			Value:       []string{"http://localhost:2379"},
			Destination: &etcdAdvertiseClientURLs,
			Sources: cli.EnvVars(
				"ETCD_ADVERTISE_CLIENT_URLS",
			),
		},
		&cli.StringSliceFlag{
			Name:        "etcd.initial-advertise-peer-urls",
			Usage:       "Initial advertise address for peer traffic. You should set to http://$(HOSTNAME).$(SERVICE_NAME)",
			Value:       []string{"http://localhost:2380"},
			Destination: &etcdInitialAdvertisePeerURLs,
			Sources: cli.EnvVars(
				"ETCD_INITIAL_ADVERTISE_PEER_URLS",
			),
		},
		&cli.StringFlag{
			Name:        "etcd.initial-cluster-state",
			Usage:       "Initial cluster state for etcd",
			Value:       "new",
			Destination: &etcdInitialClusterState,
			Sources: cli.EnvVars(
				"ETCD_INITIAL_CLUSTER_STATE",
			),
		},
		&cli.StringFlag{
			Name:        "etcd.initial-cluster",
			Usage:       "Initial cluster for etcd",
			Value:       "",
			Destination: &etcdInitialCluster,
			Sources: cli.EnvVars(
				"ETCD_INITIAL_CLUSTER",
			),
		},
		&cli.StringFlag{
			Name:        "etcd.data-dir",
			Usage:       "Path to etcd data directory",
			Value:       "data",
			Destination: &etcdDataDir,
			Sources: cli.EnvVars(
				"ETCD_DATA_DIR",
			),
		},
		&cli.StringFlag{
			Name:        "etcd.initial-cluster-token",
			Usage:       "Initial cluster token for etcd",
			Value:       "default",
			Destination: &etcdInitialClusterToken,
			Sources: cli.EnvVars(
				"ETCD_INITIAL_CLUSTER_TOKEN",
			),
		},
	},
	Action: func(ctx context.Context, _ *cli.Command) error {
		cfg := embed.NewConfig()
		cfg.Dir = etcdDataDir
		cfg.Name = etcdName
		cfg.ListenClientUrls = stringSliceToURLs(etcdListenClientURLs)
		cfg.ListenPeerUrls = stringSliceToURLs(etcdListenPeerURLs)
		cfg.AdvertiseClientUrls = stringSliceToURLs(etcdAdvertiseClientURLs)
		cfg.AdvertisePeerUrls = stringSliceToURLs(etcdInitialAdvertisePeerURLs)
		cfg.ClusterState = etcdInitialClusterState
		cfg.InitialCluster = etcdInitialCluster

		// TODO: start etcd here

		<-ctx.Done()

		slog.Info("shutting down", "reason", ctx.Err())

		return nil
	},
}

func stringSliceToURLs(s []string) []url.URL {
	var urls []url.URL
	for _, s := range s {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		urls = append(urls, *u)
	}
	return urls
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx, os.Args); err != nil {
		slog.Error("failed to run", "error", err.Error())
		os.Exit(1)
	}
}
```

We can add the following code in the `TODO` to start etcd:

```go
		slog.Info("starting etcd")
		e, err := embed.StartEtcd(cfg)
		if err != nil {
			return fmt.Errorf("failed to start etcd: %w", err)
		}
		defer e.Close()
```

Then, health check it before running additional commands:

```go
		select {
		case <-e.Server.ReadyNotify():
			slog.Info("etcd started")
		case <-time.After(60 * time.Second):
			e.Server.Stop() // trigger a shutdown
			return errors.New("timed out waiting for etcd to start")
		}
```

And at this point, we can already test etcd cluster with Docker compose.

### Testing etcd cluster with Docker compose

Here's the configuration:

```yaml
services:
  app1:
    image: golang:1.26.4
    working_dir: /app
    volumes:
      - ./:/app
      - app1-data:/data
    environment:
      ETCD_NAME: 'app1'
      ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
      ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app1:2380'
      ETCD_INITIAL_CLUSTER_STATE: 'new'
      ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
      ETCD_DATA_DIR: '/data'
      ETCD_INITIAL_CLUSTER_TOKEN: 'default'

      # We expose etcd on port 2379 for testing
      ETCD_LISTEN_CLIENT_URLS: 'http://0.0.0.0:2379'
      ETCD_ADVERTISE_CLIENT_URLS: 'http://localhost:2379'
    ports:
      - '2379:2379'
    command: ['go', 'run', './main.go']

  app2:
    image: golang:1.26.4
    working_dir: /app
    volumes:
      - ./:/app
      - app2-data:/data
    environment:
      ETCD_NAME: 'app2'
      ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
      ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app2:2380'
      ETCD_INITIAL_CLUSTER_STATE: 'new'
      ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
      ETCD_DATA_DIR: '/data'
      ETCD_INITIAL_CLUSTER_TOKEN: 'default'
    command: ['go', 'run', './main.go']

  app3:
    image: golang:1.26.4
    working_dir: /app
    volumes:
      - ./:/app
      - app3-data:/data
    environment:
      ETCD_NAME: 'app3'
      ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
      ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app3:2380'
      ETCD_INITIAL_CLUSTER_STATE: 'new'
      ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
      ETCD_DATA_DIR: '/data'
      ETCD_INITIAL_CLUSTER_TOKEN: 'default'
    command: ['go', 'run', './main.go']

volumes:
  app1-data:
  app2-data:
  app3-data:
```

Bring up:

```bash
docker compose up -d
```

Then, install `etcdctl` to make health checks, and run:

```bash
# Put foo=bar in the store
etcdctl --endpoints=http://localhost:2379 put foo bar
# Fetch foo
etcdctl --endpoints=http://localhost:2379 get foo
# Check the members
etcdctl --endpoints=http://localhost:2379 member list
```

Try some chaos engineering, make `app1` crash:

```bash
docker compose restart app1
# Check the logs
docker compose logs app1
```

You should be able to see that `app1` is able to rejoin the cluster. You can
also expose `app2` to check if the data is consistent.

At this point, you've already initialized the etcd cluster. Since we are
embedding, there is no need to actually expose the client port. To interact with
the embedded cluster, simply use `go.etcd.io/etcd/client/v3` as usual.

### Handling user authentication

Let's use JWT and a PEM-encoded RSA public key to handle user
authentication.

```go
// auth/middleware.go
package auth

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type contextKey string

const ClaimsContextKey contextKey = "jwt_claims"

// Middleware validates JWTs against a PEM public key
func Middleware(pemKey []byte) (func(http.Handler) http.Handler, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}
			tokenStr := parts[1]

			tok, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.RS256, jose.ES256})
			if err != nil {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			var claims jwt.Claims
			if err := tok.Claims(pubKey, &claims); err != nil {
				http.Error(w, "Invalid token signature or claims", http.StatusUnauthorized)
				return
			}

			expected := jwt.Expected{
				Time: time.Now(),
			}
			if err := claims.Validate(expected); err != nil {
				http.Error(w, "Token validation failed: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

```

We're injecting the claims in the context. Let's add a helper function to get the
claims from the context:

```go
func FromContext(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(jwt.Claims)
	return claims, ok
}
```

### Setting up the HTTP handler and server

Set up the handlers:

```go
// store/handler.go
package store

import (
	"embedded-etcd/auth"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type PostRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func namespaceKey(user string, key string) string {
	return fmt.Sprintf("/%s/%s", user, key)
}

func PostHandler(cli *clientv3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Fetch user from context
		claims, ok := auth.FromContext(ctx)
		if !ok {
			panic("auth middleware not called")
		}

		// Parse payload from JSON body
		var payload PostRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Serialize for etcd
		var value strings.Builder
		if err := json.NewEncoder(&value).Encode(payload.Value); err != nil {
			slog.Error("failed to serialize value", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write to etcd
		key := namespaceKey(claims.Subject, payload.Key)
		_, err := cli.Put(ctx, key, value.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetHandler(cli *clientv3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Fetch user from context
		claims, ok := auth.FromContext(ctx)
		if !ok {
			panic("auth middleware not called")
		}

		// Parse payload from query parameters
		q := r.URL.Query()
		k := q.Get("key")
		if k == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		// Read from etcd
		key := namespaceKey(claims.Subject, k)
		resp, err := cli.Get(ctx, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(resp.Kvs) == 0 {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}

		// Value is already serialized to JSON. So, return it directly.
		w.Write(resp.Kvs[0].Value)
	}
}

```

Nothing special. Notice I prefer to panic if auth middleware is not called.
Since it's a critical developer error, it's better to panic to notify the
developer.

Set up the server in the main:

```go
// main.go
var (
	// ...

	httpListenAddr string
	publicKeyPath  string
)

var app = &cli.Command{
	// ...
	Flags: []cli.Flag{
		// ...
		&cli.StringFlag{
			Name:        "http.listen-addr",
			Usage:       "Listen address for HTTP traffic.",
			Value:       ":3000",
			Destination: &httpListenAddr,
			Sources: cli.EnvVars(
				"HTTP_LISTEN_ADDR",
			),
		},
		&cli.StringFlag{
			Name:        "jwt.public-key",
			Usage:       "Path to public key for JWT authentication.",
			Required:    true,
			Destination: &publicKeyPath,
			Sources: cli.EnvVars(
				"JWT_PUBLIC_KEY",
			),
		},
	},
	Action: func(ctx context.Context, _ *cli.Command) error {
		publicKey, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read public key: %w", err)
		}

		// ...

		// Replace the <-ctx.Done() with the following:
		auth, err := auth.Middleware(publicKey)
		if err != nil {
			return fmt.Errorf("failed to create auth middleware: %w", err)
		}
		http.Handle("GET /data", auth(store.GetHandler(cli)))
		http.Handle("POST /data", auth(store.PostHandler(cli)))

		slog.Info("starting server", "addr", httpListenAddr)
		if err := http.ListenAndServe(httpListenAddr, nil); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve HTTP server: %w", err)
		}

		// ...
	}
}

// ...
```

### Testing the application

Update the `docker-compose`:

```diff
 services:
   app1:
     image: golang:1.26.4
     working_dir: /app
     volumes:
       - ./:/app
       - app1-data:/data
+      - ./keys:/keys:ro
     environment:
       ETCD_NAME: 'app1'
       ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
       ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app1:2380'
       ETCD_INITIAL_CLUSTER_STATE: 'new'
       ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
       ETCD_DATA_DIR: '/data'
       ETCD_INITIAL_CLUSTER_TOKEN: 'default'

       # We expose etcd on port 2379 for testing
       ETCD_LISTEN_CLIENT_URLS: 'http://0.0.0.0:2379'
       ETCD_ADVERTISE_CLIENT_URLS: 'http://localhost:2379'
+      HTTP_LISTEN_ADDR: ':3000'
+      JWT_PUBLIC_KEY: '/keys/public.pem'
     ports:
       - '2379:2379'
+      - '3000:3000'
     command: ['go', 'run', './main.go']

   app2:
     image: golang:1.26.4
     working_dir: /app
     volumes:
       - ./:/app
       - app2-data:/data
+      - ./keys:/keys:ro
     environment:
       ETCD_NAME: 'app2'
       ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
       ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app2:2380'
       ETCD_INITIAL_CLUSTER_STATE: 'new'
       ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
       ETCD_DATA_DIR: '/data'
       ETCD_INITIAL_CLUSTER_TOKEN: 'default'
+      JWT_PUBLIC_KEY: '/keys/public.pem'
     command: ['go', 'run', './main.go']

   app3:
     image: golang:1.26.4
     working_dir: /app
     volumes:
       - ./:/app
       - app3-data:/data
+      - ./keys:/keys:ro
     environment:
       ETCD_NAME: 'app3'
       ETCD_LISTEN_PEER_URLS: 'http://0.0.0.0:2380'
       ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://app3:2380'
       ETCD_INITIAL_CLUSTER_STATE: 'new'
       ETCD_INITIAL_CLUSTER: 'app1=http://app1:2380,app2=http://app2:2380,app3=http://app3:2380'
       ETCD_DATA_DIR: '/data'
       ETCD_INITIAL_CLUSTER_TOKEN: 'default'
+      JWT_PUBLIC_KEY: '/keys/public.pem'
     command: ['go', 'run', './main.go']

 volumes:
   app1-data:
   app2-data:
   app3-data:
```

Generate the public key:

```bash
mkdir -p ./keys
openssl genrsa -out ./keys/priv.pem 2048
openssl rsa -in ./keys/priv.pem -pubout -out ./keys/public.pem
```

Re-bring up:

```bash
docker compose up -d
```

Test without authentication:

```bash
curl 'http://localhost:3000/data?key=test'
# Authorization header missing
```

Let's create a JWT:

```bash
base64url() {
  openssl base64 -A | tr '+/' '-_' | tr -d '='
}

SUBJECT="my-user"
PRIVATE_KEY="./keys/priv.pem"

header='{"alg":"RS256","typ":"JWT"}'
exp=$(( $(date +%s) + 300 ))
payload=$(printf '{"sub":"%s","exp":%s}' "$SUBJECT" "$exp")

# Encode
header_b64=$(printf '%s' "$header" | base64url)
payload_b64=$(printf '%s' "$payload" | base64url)
unsigned="${header_b64}.${payload_b64}"

# Sign
signature=$(printf '%s' "$unsigned" | openssl dgst -sha256 -sign "$PRIVATE_KEY" | base64url)

token="${unsigned}.${signature}"
```

Now send the request with the JWT:

```bash
curl -H "Authorization: Bearer $token" 'http://localhost:3000/data?key=test'
# key not found
```

We can confirm the authentication works, now let's create a key:

```bash
curl -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{"key": "my-key", "value": "my-value"}' \
  'http://localhost:3000/data'
```

Now let's fetch the key:

```bash
curl -H "Authorization: Bearer $token" 'http://localhost:3000/data?key=my-key'
# "my-value"
```

It works!

Let's use `etcdctl` to inspect the store:

```bash
etcdctl --endpoints=http://localhost:2379 get --prefix '/my-user'
# /my-user/my-key
# "my-value"
```

Pretty cool, huh?

## A small drawback

Etcd member management is not dynamic. The only way to add a new member is to
use `etcdctl member add`. You can also do this with code, but, what I mean, is
that there is no auto-discovery of new members.

You **need** to install `etcdctl` alongside your program to be able to manage
the embedded etcd cluster, in case of disaster.

## Conclusion

Thanks to etcd, we were able to create a distributed stateful application with
high availability.

There is many other stores you can embed, like NATS JetStream for example.
Pretty useful when creating a lightweight package, huh?

This is not limited to Go, you can also use Infinispan for Java. But the main
requirement is that the developed application must use the same programming
language as the store's programming language.

Anyway, I hope this article has been useful for you.
