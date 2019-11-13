# Some notes for Kubernetes starters

This md is to introduce some basic commands for using a kubernetes cluster.

## Kubectl

This is the main tool to use the cluster.

### Get

`kubectl get [resource kind] [-n namespace]`

this will show the resources of provided kind in provided namespace. Note that `[-n namespace]` is optional, the default is `default`, and `-A` means to show resources of all namespaces.

For this assignment, `kube-system` is the most likely used namespace

examples: `kubectl get pods -n kube-system`, `kubectl get pods,services,deployments -A`

### Apply

`kubectl apply -f [filename]`

This will apply the yaml file to cluster.

examples: `kubectl apply -f apiserver_deployment.yaml`

### More Details

The above two is the mostly used command, as nearly all things can be done with a yaml file containing proper configs on a proper initialized cluster.

For more details this web page is useful.
<https://blog.csdn.net/bbwangj/article/details/80814568>

## Yaml files

Yaml files is the config file for kubernetes clusters. Here is a example: 

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-scheduler
    tier: control-plane
  name: my-kube-scheduler
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-scheduler
    - --bind-address=127.0.0.1
    - --kubeconfig=/my-scheduler/config.yaml
    image: k8s.gcr.io/kube-scheduler:v1.15.0
    imagePullPolicy: IfNotPresent
    livenessProbe:
      failureThreshold: 8
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10251
        scheme: HTTP
      initialDelaySeconds: 15
      timeoutSeconds: 15
    name: kube-scheduler
    resources:
      requests:
        cpu: 100m
    volumeMounts:
    - mountPath: /my-scheduler
      name: my-scheduler-config
  hostNetwork: true
  priorityClassName: system-cluster-critical
  volumes:
  - name: my-scheduler-config
    configMap:
      name: my-scheduler-config
status: {}
```

This is a yaml for a scheduler pod. As it shows, nearly all things can be configed in this file. 

Some key points: 
1. `containers` part: it tells the cluster which docker iamge will be used, and how to use it: `image` part(use what) and `command` part(how to use, in another word, when manully runs docker, what command you type into the terminal).
2. `volumes` part: volumes that will be mounted to dockers. It may be a path or configmap or something else.
3. `name` and `namespace`: How you identify the resource.

For more details you can refer to this page: <https://www.cnblogs.com/lgeng/p/11053063.html>
