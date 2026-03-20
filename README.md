# rbg-api

`rbg-api` is the standalone API and client library for [RoleBasedGroup (RBG)](https://sigs.k8s.io/rbgs) — a Kubernetes operator for managing role-based LLM inference workloads (e.g., prefill/decode disaggregation).

This repository is intended to be imported by **upper-layer business systems** that need to interact with RBG resources via the Go client, without depending on the full operator codebase.

## Module

```
sigs.k8s.io/rbgs/api
```

## Project Structure

```
rbg-api/
├── apis/                            # API type definitions
│   └── workloads/
│       ├── constants/               # Shared constants (labels, annotations, enums)
│       ├── v1alpha1/                # workloads.x-k8s.io/v1alpha1 API types
│       │   ├── doc.go
│       │   ├── groupversion_info.go
│       │   └── ...types.go
│       └── v1alpha2/                # workloads.x-k8s.io/v1alpha2 API types (storage version)
│           ├── doc.go
│           ├── groupversion_info.go
│           ├── rolebasedgroup_types.go
│           ├── rolebasedgroupset_types.go
│           ├── rolebasedgroupscalingadapter_types.go
│           ├── roleinstance_types.go
│           └── roleinstanceset_types.go
├── client-go/                       # Generated Go client libraries
│   ├── clientset/versioned/         # Typed clientset
│   │   ├── clientset.go
│   │   ├── scheme/                  # Scheme registration
│   │   └── typed/workloads/
│   │       ├── v1alpha1/            # v1alpha1 typed clients
│   │       └── v1alpha2/            # v1alpha2 typed clients
│   ├── informers/externalversions/  # SharedInformerFactory
│   │   ├── factory.go
│   │   ├── generic.go
│   │   ├── internalinterfaces/
│   │   └── workloads/
│   │       ├── v1alpha1/
│   │       └── v1alpha2/
│   └── listers/workloads/           # Listers
│       ├── v1alpha1/
│       └── v1alpha2/
├── config/crd/bases/                # Generated CRD YAML manifests
├── hack/
│   ├── boilerplate.go.txt           # License header for code generation
│   └── update-codegen.sh            # Script to regenerate client-go code
├── Makefile
└── go.mod
```

## Key Resources

### API Group: `workloads.x-k8s.io`

| Kind | Version | Short Name | Description |
|------|---------|------------|-------------|
| `RoleBasedGroup` | v1alpha2 | `rbg` | Core resource: defines multiple named roles (e.g., prefill, decode), each with replicas and pod templates |
| `RoleBasedGroupSet` | v1alpha2 | `rbgs` | Manages a set of identical RoleBasedGroups |
| `RoleBasedGroupScalingAdapter` | v1alpha2 | `rbgsa` | HPA-compatible scaling adapter for a specific role in an RBG |
| `RoleInstance` | v1alpha2 | `rins` | A single replica unit composed of one or more pods |
| `RoleInstanceSet` | v1alpha2 | `ris` | Manages a set of RoleInstances for a role |

## Installation

```bash
go get sigs.k8s.io/rbgs/api@latest
```

## Usage

### Create a clientset

```go
import (
    "k8s.io/client-go/tools/clientcmd"
    rbgclient "sigs.k8s.io/rbgs/api/client-go/clientset/versioned"
)

config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
    panic(err)
}

clientset, err := rbgclient.NewForConfig(config)
if err != nil {
    panic(err)
}
```

### List RoleBasedGroups

```go
import (
    "context"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

rbgList, err := clientset.WorkloadsV1alpha2().RoleBasedGroups("default").List(
    context.Background(),
    metav1.ListOptions{},
)
```

### Use Informers

```go
import (
    "time"
    "k8s.io/client-go/tools/cache"
    rbginformers "sigs.k8s.io/rbgs/api/client-go/informers/externalversions"
)

factory := rbginformers.NewSharedInformerFactory(clientset, 30*time.Second)
rbgInformer := factory.Workloads().V1alpha2().RoleBasedGroups().Informer()

rbgInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
        // handle add
    },
    UpdateFunc: func(oldObj, newObj interface{}) {
        // handle update
    },
    DeleteFunc: func(obj interface{}) {
        // handle delete
    },
})

stopCh := make(chan struct{})
factory.Start(stopCh)
factory.WaitForCacheSync(stopCh)
```

### Use Listers

```go
import (
    "k8s.io/apimachinery/pkg/labels"
    rbglistersv2 "sigs.k8s.io/rbgs/api/client-go/listers/workloads/v1alpha2"
    "k8s.io/client-go/tools/cache"
)

lister := rbglistersv2.NewRoleBasedGroupLister(indexer)
rbgs, err := lister.RoleBasedGroups("default").List(labels.Everything())
```

## Code Generation

The `client-go` packages are generated from the API types. To regenerate after modifying API types:

```bash
# Install required tools
make install-tools

# Regenerate deepcopy methods and CRD manifests
make generate

# Regenerate client-go (clientset, listers, informers)
make generate-clients
```

## API Versioning

| Version | Status | Notes |
|---------|--------|-------|
| `v1alpha1` | Deprecated | Legacy API, use v1alpha2 for new code |
| `v1alpha2` | Active | Current storage version, recommended |

## License

Apache License 2.0 — see [LICENSE](../LICENSE).
