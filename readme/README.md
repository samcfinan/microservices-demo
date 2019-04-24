# README

## Overview

This repo contains a basic Kubernetes-based microservices demo primarily written in Golang. gRPC transports are used exclusively between services and defined gRPC requests can be received by the external load balancer. Please read the sections below to understand the technologies used, structure, deployment, and suggested development workflow.

Skip to the [Setup for Dev](#setup-for-dev) section if you don't like to read.

## Containers

Containers are lightweight self-contained environment that package only what an application or service needs to run. They are extremely light, and unlike VMs, share the host kernel and memory address space. By packaging the entire app's dependencies, they are highly portable and become hardware and OS-agnostic.

[Docker](https://www.docker.com/why-docker) is by far the most popular container platform.

## Kubernetes (k8s)

All containerized services (built as Docker images), including the load balancer, are orchestrated by Kubernetes. Kubernetes groups deployments (individual microservices) into Pods and manages container restarts, healing, load balancing, and replication automatically, and can scale from single-node clusters to the datacenter-scale. This is handled declaratively via manifests (YAML-formatted descriptions of deployments and exposed services).

Kubernetes manifests declaring deployments, services, and exposed ports are found in the [kubernetes-manifests](/kubernetes-manifests) directory.

Kubernetes relies on Docker to handle containerization. A `Dockerfile` is found in each service's repo which should declare its base image, build process, exposed ports, and entrypoint.

[Read More](https://kubernetes.io/)

## Skaffold

Skaffold is a command line tool that facilitates easier local development for Kubernetes. It is based on a single declarative YAML file stored in the root of your project. Note that Skaffold does not replace Docker, Kubernetes, or Kubectl, but manages image configuration and hot reloading during your development process.

[skaffold.dev](https://skaffold.dev/)\
[github.com/GoogleContainerTools/skaffold](https://github.com/GoogleContainerTools/skaffold)

## gRPC

This project's internal transports are handled via [gRPC](https://grpc.io/). gRPC defines services using Protocol Buffers, standardizing request-response contracts between languages. Protobuf's binary serialization reduces message size but up to 10x in normal applications.

Protobufs, written in proto3, are found in the [pb](/pb) directory. These definitions must be compiled to be read by the service's language. Each service's directory contains a `genproto.sh` bash script to compile to that language's format.

#### Note: libprotoc 3.7.0 must be installed to compile to output languages.

## Setup For Dev

Clone repo into your Gopath. For example: `$GOPATH/src/github.com/samcfinan/microservices-demo`.

Install the following dependencies:

#### Note: this guide assumes you are using Linux/MacOS. Windows containerization is somewhat different and your mileage may vary.

1. [Skaffold](https://skaffold.dev/docs/getting-started/#installing-skaffold): Declarative Kubernetes setup and development environment.
2. [Docker](https://docs.docker.com/install/): Containerization tool.
3. [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/): Kubernetes CLI admin.
4. [Minikube](https://kubernetes.io/docs/setup/minikube/#installation): Single-node Kubernetes cluster ideal for local development.
5. [Protocol Buffers](https://github.com/protocolbuffers/protobuf): Binary serialization for gRPC. Install a pre-built binary from the [releases](https://github.com/protocolbuffers/protobuf/releases) page.

You will also require the language binary for any language you plan to write a service in, of course.

Run the demo:

1. Download and vendor dependencies for Go services. `cd` into each Go service and run `go mod vendor`. Note that your `GO111MODULE` environment variable must be set to `on`.
2. Start minikube with sufficient resources. `minikube start --memory=4000 --cpus=3 --disk-size=32g`
3. Run `kubectl get nodes` to verify you're connected to “Kubernetes on Docker”.
4. Run `skaffold run` (first time will be slow, it can take ~5-10 minutes). This will build and deploy the application. If you need to rebuild the images automatically as you refactor the code, run `skaffold dev` command.
5. Run `kubectl get services` or `minikube dashboard` to find your service's exposed IP.

Containers will rebuild and k8s will port forward to the new container on change detection. Please note that this may take a few seconds.

## Setup for Prod

K8s is ideal for a cloud-native deployment but can be used on-premise as well. 

Skaffold relies on Kubectl to build and configure the deployments. As such, Kubectl needs to connect to the Kubernetes API on the cluster's master node.

[https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/)

From there, Skaffold will operate as though Kubernetes is running locally via Minikube. Note that the K8s host must expose port 443 and be set up with an SSL certificate.

Verify that your configuration is accurate by running `kubectl config view` to print the config and `kubectl get nodes` to list the available nodes in the cluster.
