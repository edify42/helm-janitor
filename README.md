# helm-janitor

_helm-janitor_ is an API interface to clean up releases in your k8s cluster.

It can also scans for [helm3](https://helm.sh/blog/helm-3-released/) releases
and runs a `helm uninstall|delete` against releases that have a `janitor-ttl`
annotation which has expired on the release.

This first cut of code is initially intended to run in the [Lendi](https://www.lendi.com.au)
k8s clusters to clean up helm releases in our development environment.

**Our setup**
- EKS (version >= 1.19)
- Helm releases stored as k8s secrets
- Bitbucket webhooks which fire when PRs are merged.

## usage as tool

Support 2 modes of running

```bash
./helm-janitor [command] [options]

[command]
delete <release>
scan

[delete-options]
--namespace <namespace>

[scan-options]:
--namespace <namespace>
--all-namespaces
--include-namespace <expression match>
--exclude-namespace <expression match>
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

### k8s container

Other k8s native option is to run this as a k8s cronjob that can remove the helm release.

#### permission requirements

Needs an RBAC cluster role binding which has the right amount of cluster
permissions to remove a helm release. 

## future work

- Support SQL backends?
- Support configmap backed helm releases?