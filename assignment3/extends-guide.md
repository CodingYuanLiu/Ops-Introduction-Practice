# Extends-guide
> a simple guide by Qingyuan for building a simple scheduler extension.

The complete code is at ./ops/extends

## Step 1: Run a simple server

The extension act like an REST service, which receives a json request from the kubenetes scheduler and send a json response back to it. So, the first step we need to do is to run a simple server written by golang.

```go
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)

	log.Fatal(http.ListenAndServe(":8888", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to sample-scheduler-extender!\n")
}

//Filter() and Prioritize() will be shown later.
...
```

The `httprouter` can be used to quickly set up a server. Just run the server by `go run main.go` (after initialize `go module`, of course), and you can get the response at `localhost:8888`

```bash
ubuntu@kubernetes:~/ops$ curl localhost:8888
Welcome to sample-scheduler-extender!
```

## Step 2: write a simple algorithm

The `Prioritize()` is written to score the node randomly. 

```go
//router.go
func Prioritize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		log.Println(err)
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = prioritize(extenderArgs)
	}

	if response, err := json.Marshal(hostPriorityList); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

//Prioritize.go
const (
	// lucky priority gives a random [0, schedulerapi.MaxPriority] score
	// currently schedulerapi.MaxPriority is 10
	luckyPrioMsg = "pod %v/%v is lucky to get score %v\n"
)

// it's webhooked to pkg/scheduler/core/generic_scheduler.go#PrioritizeNodes()
// you can't see existing scores calculated so far by default scheduler
// instead, scores output by this function will be added back to default scheduler
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		score := rand.Intn(schedulerapi.MaxPriority + 1)
		log.Printf(luckyPrioMsg, pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
```

A simple filter algorithm to filter the node randomly is also easy to implement. The code can be seen at github.

#### Cautions

The GOPROXY updated its content with golang 1.13, as every module needs a version number. However, if you don't assign a version, a v0.0.0 number will automatically assigned, which may cause that the GOPROXY can not find the resources and replies "410 Gone". As a result, we should replace all the module's name with a **latest** version number in `go.mod` file:

```go mod
replace k8s.io/cli-runtime => k8s.io/cli-runtime latest
```

After running, this entry will automatically find the latest version number and replace the *latest* with it.

```go mod
replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191115221445-ec04ad4dbd24
```



## Step 3: build the image with the dockerfile

In order to connect the kube-scheduler, the extension is expected to run in the same pod with kube-scheduler, so that it needs to run in a docker container. Build the image with the dockerfile:

```dockerfile
FROM golang:1.12
# COPY the souce code of the project
COPY ./*.go /app/
# COPY the go module file of the project
COPY ./go.* /app/
ENV GOPROXY=http://goproxy.io
EXPOSE 8888
WORKDIR /app

ENTRYPOINT ["go","run","main.go"]
```

then: simply `go build -t <imagename> ."  is okay.

## Step 4: apply the scheduler with the extension

Modify the scheduler's yaml to run the container in the same pod with extensions.

my-scheduler.yaml:

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
    - --kubeconfig=/etc/kubernetes/scheduler.conf
    - --leader-elect=false
    - --config=/my-scheduler/config.yaml
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
    - mountPath: /etc/kubernetes/scheduler.conf
      name: kubeconfig
      readOnly: true
    - mountPath: /my-scheduler
      name: my-scheduler-config
  - image: lqyuan980413/k8s_scheduler_extends
    imagePullPolicy: IfNotPresent
    name: myextender  
  hostNetwork: true
  priorityClassName: system-cluster-critical
  volumes:
  - hostPath:
      path: /etc/kubernetes/scheduler.conf
      type: FileOrCreate
    name: kubeconfig
  - name: my-scheduler-config
    configMap:
      name: my-scheduler-config
status: {}

```

Of course, modify the scheduler-configuration yaml too:

my-scheduler-config:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-scheduler-policy
  namespace: kube-system
data:
 policy.cfg : |
  {
    "kind" : "Policy",
    "apiVersion" : "v1",
    "predicates" : [
    ],
    "priorities" : [
    ],
    "extenders":[{
        "urlPrefix": "http://localhost:8888/",
        "filterVerb": "filter",
        "prioritizeVerb": "prioritize",
        "weight":1,
        "enableHttps":false
    }],
    "hardPodAffinitySymmetricWeight" : 10
  }

```

Then apply the kube-scheduler and start a new pod. You can see the pod can be scheduled to a new node by your scheduler.

```bash
# delete the former scheduler
kubectl delete pod my-kube-scheduler -n kube-system
# apply the new scheduler
kubectl apply -f my-scheduler-config.yaml
kubectl apply -f my-scheduler.yaml
# test
kubectl apply -f testA.yaml
```

