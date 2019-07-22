# Bugle


This allows another pods log file to be viewed by HTTP. This is essential for development teams that need to review logs only available in k8s but do not have access to kubernetes.

The container is published to docker hub.


## Instructions for use

If you want to deploy this app you only require the helm  chart. Everything else is available in the container on Docker Hub.

In the values the following parameters will need to be set

```yaml
params:
  namespace: default #Namespace is the target pod in
  podname: kubernetes-bootcamp-6bf84cb898-hsg6w #Name of the target pod
  username: unset #Credentials to secure webpage showing the log. If either of these are set to 'unset' then there is no challenge.
  password: unset
```

Ingress should also be configured in the values file. If TLS is requried (which it should) it must be enabled here. I have not put in a default SSL config as every person has a slightly different config


## Dir

* app  - The go Application
* container - The DockerFile
* helm - The helm chart


## Todo
1. Add facility for Oauth or other identiy provider
2. Filter String
3. Auto Refresh
