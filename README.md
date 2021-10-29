# helm-janitor
[![Go Report Card](https://goreportcard.com/badge/github.com/lendi-au/helm-janitor)](https://goreportcard.com/report/github.com/lendi-au/helm-janitor)
![Helm+Janitor](art/helm+janitor.png)

_helm-janitor_ is an API interface to clean up releases in your k8s cluster.

It can also scans for [helm3](https://helm.sh/blog/helm-3-released/) releases
and runs a `helm uninstall|delete` against releases that have a `janitor-ttl`
annotation which has expired on the release.

This first cut of code is initially intended to run in the [Lendi](https://www.lendi.com.au)
k8s clusters to clean up helm releases in our development environment.

The goal of this project is very much to support the Lendi development
workflow and has been built around the infrastructure here.

**Our setup**
- EKS (version >= 1.19)
- Helm releases stored as k8s secrets
- Bitbucket webhooks which fire when PRs are merged.

## usage as tool

Support 2 modes of running (delete | scan)

```bash
./helm-janitor [command] [options]

[command]
delete <selector> # helm-janitor=true
scan <selector> # BRANCH=feat/test-something,REPOSITORY=cool-repo

[options]
--namespace <namespace>
--all-namespaces (default)
--include-namespace <expression match>
--exclude-namespace <expression match>
```

`scan` will scan all namespaces if no specific namespace or expressions are
used.

### delete use-case

When we run our k8s deployments via [Helm](https://helm.sh), we also tag and
label the helm releases (secrets) with our tooling using the repository and
the branch that the release was deployed from.

Teams have a webhook configured on there repo which fires when a PR merged /
branch is closed and [Stack Janitor](https://github.com/lendi-au/stackjanitor)
will clean up the left over running containers.

### scan use-case

Like the [kube-janitor project](https://codeberg.org/hjacobs/kube-janitor),
we wish to expire helm releases that exceed the `ttl` value. During our
`helm instal...` step, we tag the release secret afterwards with a
`helm-janitor: true` label if we wish to clean up the release. We then read the
release values `janitorAnnotations` config to check the `janitor/ttl` or
`janitor/expires` values and checks against the creationTime to see if we
should delete the release.

```yaml
janitorAnnotations:
  janitor/expires: "2021-07-03T07:06:45Z"
```

## k8s mgmt use-cases/ecosystem

- Run this when teams want to manually clean up release via slack-ops
- Integrate with CI/CD systems via webhook call
- Custom cleanup schedule on certain environments

## deployment models

### AWS lambda

We can use an AWS lambda to periodically run the app (via serverless) in our
AWS environment.

#### permission requirements

Need an IAM role that maps to a k8s RBAC cluster role which has enough
permissions to clean up the helm release.

Use the `ROLE_ARN` environment variable if the lambda needs to assume a role to
access the k8s cluster.

We use map the AWS IAM role to a cluster user which has the sufficient RBAC
permission to remove a release.

### k8s container

Other k8s native option is to run this as a k8s cronjob that can remove the
helm release.

#### permission requirements

Needs an RBAC cluster role binding which has the right amount of cluster
permissions to remove a helm release.

## contributing

We may reject PRs that break compatibility with our k8s setup.

## future work

- Support SQL backends?
- Support configmap backed helm releases?