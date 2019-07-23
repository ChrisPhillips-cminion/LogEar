# LogEar


For full introduction of LogEar please go to  (https://chrisphillips-cminion.github.io/kubernetes/2019/07/23/LogEar.html)[https://chrisphillips-cminion.github.io/kubernetes/2019/07/23/LogEar.html]

LogEar is a container that runs in kubernetes and allows access to the log of a single other container via a webpage. Multiple instances can be run if multiple container logs are required.

LogEar works by calling the kubernetes API to access the logs on demand. The user initiates a request via HTTPs and this triggers the LogEar Go application inside of its container. The log is then returned to the HTTP requestor.


The container is published to docker hub and so can be downloaded directly from there instead of building from this source.


## Instructions for use

If you want to deploy this app you only require the helm chart. Everything else is available in the container on Docker Hub. (http://hub.docker.com/u/cminion/LogEar)[http://hub.docker.com/u/cminion/LogEar]

In the values file for the helm chart the following parameters will need to be set

```yaml
params:
  namespace: default #Namespace is the target pod in
  podname: kubernetes-bootcamp-6bf84cb898-hsg6w #Name of the target pod
  username: unset #Credentials to secure webpage showing the log. If either of these are set to 'unset' then there is no challenge.
  password: unset
```

Ingress should also be configured in the values file. If TLS is requried (which it should) it must be enabled here. I have not put in a default SSL config as every organization has a slightly different config


## Repo Directory Structure

* app  - The go Application
* container - The DockerFile
* helm - The helm chart


## Todo
1. Add facility for Oauth or other identiy provider
2. Filter String
3. Auto Refresh
