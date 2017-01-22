# Marathon Resource


[![Build Status](https://travis-ci.org/ckaznocha/marathon-resource.svg?branch=master)](https://travis-ci.org/ckaznocha/marathon-resource)
[![Coverage Status](https://coveralls.io/repos/github/ckaznocha/marathon-resource/badge.svg?branch=master)](https://coveralls.io/github/ckaznocha/marathon-resource?branch=master)
[![Code Climate](https://codeclimate.com/github/ckaznocha/marathon-resource/badges/gpa.svg)](https://codeclimate.com/github/ckaznocha/marathon-resource)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://ckaznocha.mit-license.org)
[![Go Report Card](https://goreportcard.com/badge/ckaznocha/marathon-resource)](https://goreportcard.com/report/ckaznocha/marathon-resource)
[![Docker Pulls](https://img.shields.io/docker/pulls/ckaznocha/marathon-resource.svg?maxAge=2592000)](https://hub.docker.com/r/ckaznocha/marathon-resource/)

A [Concourse](https://concourse.ci/) resource to deploy applications to [Marathon](https://mesosphere.github.io/marathon/).

## Source Configuration

*   `app_id`: *Required.* The name of your app in Marathon.

*   `uri`: *Required.* The URI of the Marathon instance you wish to deploy to.

*   `basic_auth`: *Optional.* Use if you are using HTTP Basic Auth to protect your Marathon instance. Takes `user_name` and `password`

*   `api_token`: *Optional.* Use if you are using DC/OS and need to set an HTTP API token.

## Behavior

### `check`: Extract versions of an app from Marathon.

Returns a list of any versions greater than or equal the last know version of the app defined by `app_id`.

### `in`: Fetch data about the current version of an app.

Returns JSON description of the current running version of the app.

#### Parameters

*None.*


### `out`: Deploy an app to Marathon.

Given a JSON file specified by `app_json`, post it to Marathon to deploy the app. The resource will cancel the deployment if its not successful after `time_out`.

#### Parameters

*   `app_json`: *Required.* Path to the JSON file describing your marathon app. For more information about the format see [the Marathon docs](https://mesosphere.github.io/marathon/docs/application-basics.html).

*   `time_out`: *Required.* How long, in seconds, to wait for Marathon to deploy the app. Timed out deployments will roll back and fail the job.

*   `replacements`: *Optional.* A `name`/`value` list of templated strings in the app.json to replace during the deploy. Useful for things such as passwords or urls that change.

*   `restart_if_no_update`: *Optional.* If Marathon doesn't detect any change in your app.json it won't deploy a new version. Setting this to `true` will restart an existing app causing a new version. Default is `false`.

## Example Configuration

### Resource type

``` yaml
- name: marathon
  type: docker-image
  source:
    repository: ckaznocha/marathon-resource
```

### Resource

``` yaml
- name: marathon_app
  type: marathon
  source:
    app_id: my_app
    uri: http://my-marathon.com/
    basic_auth:
      user_name: my_name
      password: {{ marathon_password }}
```

### Plan

``` yaml
- get: marathon_app
```

``` yaml
- put: marathon_app
  params:
    app_json: path/to/app.json
    time_out: 10
    replacements:
    - name: db_password
      value: {{ db_password }}
    - name: db_url
      value: {{ db_url }}
```

## Contributing

See the `CONTRIBUTING` file.

## License
See `LICENSE` file
