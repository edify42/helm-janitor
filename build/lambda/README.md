# lambda

Deploy our helm-janitor application as a lambda function that periodically runs
and can be invoked via HTTP method to clean the releases

## schedule

First cut of the schedule code will simply scan releases with the label
`helm-janitor: true` and look for the TTL/expires value in the values file
that was deployed with the release.


## deploy

STILL WIP - relying on Lendi specific tooling.

```bash

liam exec -e development -r deployer -- serverless deploy

```