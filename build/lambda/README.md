# lambda

Deploy our helm-janitor application as a lambda function that periodically runs
and can be invoked via HTTP method to clean the releases

## schedule

First cut of the schedule code will simply scan releases with the label
`helm-janitor: true` and look for the TTL/expires value in the values file
that was deployed with the release.

## api

Eventually we'd like an API gateway endpoint to expose the function for
further integration points (e.g. slack webhook)

## custom

The function can also be manually invoked by say, another lambda to perform
the janitor duties.

## deployment to AWS

STILL WIP - relying on Lendi specific tooling.

```bash

liam exec -e development -r deployer -- serverless deploy

```