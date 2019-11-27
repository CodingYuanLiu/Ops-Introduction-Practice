# Evaluation of cluster scheduler

We have implemented a basic scheduler by checking conditions and random priorities.

The scheduler will schedule pods by the rules below:

> 1. Check whether the pod name length is smaller than each node's index + 17.
> 2. If true, all nodes which fulfill conditions will generate a random priority, the pod will be settled in the highest-priority node.
> ( The node's index are 1, 2, 3. )

### Pod 1 (Name: simple-tester)
The pod are fit for all nodes, because its name length is 13.

Start the pod then type `kubectl logs my-kube-scheduler myextender -n kube-system`, the log is shown below:
```
2019/11/27 14:44:08 [kubeslave1] pod simple-tester/default length is 13, node length is 18, fit
2019/11/27 14:44:08 [kubeslave2] pod simple-tester/default length is 13, node length is 19, fit
2019/11/27 14:44:08 [kubeslave3] pod simple-tester/default length is 13, node length is 20, fit
2019/11/27 14:44:08 [kubeslave1] pod simple-tester/default is lucky to get score 8
2019/11/27 14:44:08 [kubeslave2] pod simple-tester/default is lucky to get score 7
2019/11/27 14:44:08 [kubeslave3] pod simple-tester/default is lucky to get score 2
```

`kubeslave1` rolled the highest number, so the pod has been allocated to `kubeslave1`. Log:
```
Events:
  Type    Reason     Age   From                 Message
  ----    ------     ----  ----                 -------
  Normal  Scheduled  98s   my-kube-scheduler    Successfully assigned default/simple-tester to kubeslave1
  Normal  Pulling    97s   kubelet, kubeslave1  Pulling image "nginx"
  Normal  Pulled     87s   kubelet, kubeslave1  Successfully pulled image "nginx"
  Normal  Created    86s   kubelet, kubeslave1  Created container podtest-scheduler
  Normal  Started    86s   kubelet, kubeslave1  Started container podtest-scheduler
```

### Pod 2 (Name: a-bit-longer-tester)
The pod are fit for node 3, because its name length is 19.

No roll needed, because there are only one node can contain the pod. Log about this:
```
2019/11/27 14:48:35 [kubeslave1] pod a-bit-longer-tester/default length is 19,  node length is 18, unfit
2019/11/27 14:48:35 [kubeslave2] pod a-bit-longer-tester/default length is 19,  node length is 19, unfit
2019/11/27 14:48:35 [kubeslave3] pod a-bit-longer-tester/default length is 19, node length is 20, fit
```

```
Events:
  Type    Reason     Age   From                 Message
  ----    ------     ----  ----                 -------
  Normal  Scheduled  119s  my-kube-scheduler    Successfully assigned default/a-bit-longer-tester to kubeslave3
  Normal  Pulling    118s  kubelet, kubeslave3  Pulling image "nginx"
  Normal  Pulled     112s  kubelet, kubeslave3  Successfully pulled image "nginx"
  Normal  Created    112s  kubelet, kubeslave3  Created container podtest-scheduler
  Normal  Started    112s  kubelet, kubeslave3  Started container podtest-scheduler

```

### Pod 3 (Name: a-looooooooooooog-tester)
The pod are unfit for all nodes, because its name length are tooooooooooo long. (24)
Log:
```
2019/11/27 14:52:11 [kubeslave1] pod a-looooooooooooog-tester/default length is 24,  node length is 18, unfit
2019/11/27 14:52:11 [kubeslave2] pod a-looooooooooooog-tester/default length is 24,  node length is 19, unfit
2019/11/27 14:52:11 [kubeslave3] pod a-looooooooooooog-tester/default length is 24,  node length is 20, unfit
```

Then check this tester by `kubectl describe pod a-looooooooooooog-tester`:
```
Events:
  Type     Reason            Age               From               Message
  ----     ------            ----              ----               -------
  Warning  FailedScheduling  5s (x2 over 79s)  my-kube-scheduler  0/4 nodes are available: 1 node(s) had taints that the pod didn't tolerate, 3 .
```