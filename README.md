# Turbine CLI

An opinionated tool for creating, developing, and deploying projects

## Commands

This API is the goal

### `setup`

Install dependancies, prepare environmennt for local development. Meant to be run on a developer's machine.

### `init ${PROJECT_NAME} --template ${template_name} --api-key ${API_KEY}`

Creates new project with database, git repo, deployment pipelines. Everything to get a service running asap.

### `dev`

Runs the project on local k8s cluster

### `run -n ${SCRIPT_NAME}`

Runs any scripts defined in `turbine.dhall`

### `deploy ${ENVIRONMENT_NAME} --arg="something"`

Deploys a given environment. Args can be passed in for deployments that need extra info. e.g. `turbine deploy feature --subdomain="pr-7"`

## Config file

`project name` - a unique name for the project
