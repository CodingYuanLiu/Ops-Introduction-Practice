# Basic Operations on Kubernetes

This md is to introduce some basic commands for using a kubernetes cluster.

## Kubectl

This is the main tool to use the cluster.

## Commands

### Get

`kubectl get [resource kind] [-namespace]`
this will show the resources of provided kind in provided namespace. Note that `[-namespace]` is optional, the default is `default`, and `-A` means to show resources of all namespaces.
examples: `kubectl get pods -Micro`, `kubectl get pods,services,deployments -A`

### Apply

`kubectl apply -f [filename]`
This will apply the yaml file to cluster.

## More Details

The above two is the mostly used command, as nearly all things can be done with a yaml file containing proper configs on a proper initialized cluster.

For more details this web page is useful.
<https://blog.csdn.net/bbwangj/article/details/80814568>
