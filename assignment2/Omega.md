# Paper Reading, Omega: flexible, scalable schedulers for large computer clusters

## Introduction

### What is Omega

Omega is a flexible, scalable scheduler for large computer cluster, but what is ***scheduler for cluster***, or say, ***cluster scheduler***

#### What is Cluster Scheduler

First, let us make a formal definition to ***Cluster***
> When a group of loosely coupled computers work together so that they can be
viewed as if they are one computer, it is called cluster [1]

简单点来说，***Cluster***，中文`集群`就是一群计算机，但是再实际应用中这群计算机被视作了单个计算机。

Then, we need to learn something about scheduler</br>
Scheduler is a procedure that repeatedly do scheduling. ***Scheduler***(中文调度器)就是在集群中进行调度的程序

> In computing, scheduling is the method by which work is assigned to resources that complete the work.[2]

`scheduling`是一种用来分配资源的方法，在集群中实现了`scheduling`的那部分程序就是scheduler啦！

#### Why Omega

##### Bottleneck of current cluster scheduler

> Increasing scale and the need for rapid response to changing requirements are hard to meet with current monolithic cluster scheduler achitectures. This restricts the rate ar which new features can be deployed, decreases efficiency and utilization, and eventually limit cluster growth

当时集群变得越来越大， 对扩展性的要求越来越高，而在集群不断变大的同时，对集群性能的需求从来没有停止过。问题是，现有的集群调度器架构并不能解决我们的问题，是的，Omega在架构上和现有的集群是不一样的，那么哪里不一样呢？为什么这些不一样呢？

##### Features in Omega

1. shared state

2. lock-free optimistic concurrency control

We will explian these two deatures in the following part of this document.

## Cluster Scheduler Design

### Metrics

1. high resource utiliaztion
2. user-supplied placement constraints
3. rapid decision making
4. various degrees of "fairness"
5. business importance

### jobs

In the paper, we use a simple way to separate jobs. We divide them into two parts.

One is ***service*** jobs that are long running jobs and take up more resource.

The other is ***batch*** jobs which quickly perform a computation and then finish.

### Design issues 

Different implementations of scheduler may solve different problems. Some sacrifice performance for accuracy, while others guarantee correctness by continuously checking. No matter what implementation it is, there are some common issues they have to dealt with. Now, we first introduce these issues. After that, we will dive into the solution of omega and other schedulers

1. Partition the scheduling work

    It is obvious that some jobs can be executed immediately, while some jobs are long running. In real world. How we treat them in our scheduler in the main part of this issue. Mainly, we have three ways: 
    
    1. Load balancing that is oblivious to workload type</br>
    2. Dedicating specialized schedulers to different parts of the workload 
2. Choice of resources
    
    We can make decisions from overall perspective or local perspective, determined by scheduler implementation. 
    
    &nbsp;&nbsp;&nbsp;&nbsp;If we decide from global perspective, we may make better decisions. However, this requires access to all of the cluster resources since scheduler must know the utilization and then allocate globally.

3. Concurrency Control

    Many schedulers are parallel, thus they have dealt with concurrency issues. 
    
    To deal with such issues, scheduler ensure that a particular resourve is only made available to one scheduler at a time. The approach mensioned is typically called pessimistic concurrency control.

    Contrasted to pessimistic approach, optimistic approach which is applied in Omega, allows a resource earned by multiple schedulers. Optimistic control detects confict, if multiple schedulers want to allocate same resource, then it undo one or more of the conficting claims.

4. Allocation granularity(分配粒度)
   
    Jobs typically contains multiple tasks, schedulers can have different policies to place one job. Typically methods are:
    all-or-nothing and incrementally acquiring

5. Cluster-wide behaviors
    Cluster-wide behaviors focus on issues span mutiple schedulers, such as achieving various types of fairness, a common agreement on precendence of work and so on.

### Omega and other schedulers

#### Monolithic Schedulers

##### What is monolithic schedulers 

Monolithic schedulers use a single centralized scheduling alogorithm for all jobs.

![monolithic-scheduler](./imgs/scheduler-arch-monolithic.png)

Not all monolithic schedulers are implemented in one scheduling logic. They can also have multiple path of logic which are called multi-path monolithic schedulers.

##### Props & Cons 

###### Props

1. Uniformed and easy to implement
2. Total access to cluster state, can make good decisions.

###### Cons

1. Not differentiate between jobs. Schedule logic should divide between jobs for better placement and lower latency.
2. Hard to scale, single instance, sindle logic. Concurrency and parallel are required.
3. Because of single procedure, there can be `head of line blocking(队头阻塞，通常出现在TCP)`. A batch job have to wait for long-running service jobs scheduling to finish, increasing latency.

### Two-level schedulers

#### What is two-level schedulers

Two-level schedulers have a single active resource manager that offers compute resources to multiple parallel, independent scheduler.

![monolithic-scheduler](./imgs/scheduler-arch-twolevel.png)

Unlike monolithic schedulers, Two-level schedulers can dynamically adjust allocation of resources to each schedulers.</br>
Two-level schedulers use a resource manager between scheduler and resource.</br>
The shchedulers behind resource manger is called ***scheduler framework***

#### Pros & Cons

##### Pros

1. Divide jobs into several types
2. separate resource placement(by resource manager) and scheduling logic(by scheduler framework)

##### Cons

1. pessimistic concurrency control is not good enough
2. unable to access to all the cluster state. For some "picky" jobs, may not do good decisions.

### Omega

Omega takes advantages of these two types of schedulers. Omega uses ***shared state*** alternatively. `shared state` allows each scheduler in Omega full access to the entire `shared state`. Omega also allows parallel and flexibility. Unlike `Two-level scheduler`, Omega uses optimistic concurrency control and there is no central resource allocator.

![monolithic-scheduler](./imgs/scheduler-arch-sharedstate.png)

* Shared state
    
    * Each scheduler in Omega maintain a private, local, frequently-updated copy of resource allocations in the cluster, which is called ***cell state***.</br>
    * Every scheduler can see the entire state of the cluster.</br>
    * Every scheduler has complete freedom to lay claim to any available cluster resources provided it has appropriate permissions and priority.</br>
    * Once a scheduler makes a placement decision
        * It updates the `shared copy` of cell state in an atomic commit. 
        * The time from state synchronization to the commit attempt is a *transaction*. Thus, at most one update commit can succeed.

* Optimistic Concurrency Controll
    * Omega schedulers operates completely in parallel and do not have to wait for jobs in other shcedulers, as a result, there is no head of line blocking(except cases where schedulers too busy)

## Evaluation

Theory shows the advantages of Omega over other schedulers, but we have to test its performance in real world workloads.

### Workloads 

* lightweight simulator

    synthetic workloads drawn from empirical workload distributions. in this workloads, we make some simplifications that this it allows us to sweep across a broad range of operating points(scheduling / resource manage) within a reasonable runtime.

* high-fidelity simulator 高保真(高仿:smile)

    high-fidelity simulator replays historic workload traces from Google production clusters and reuses much of the Google production scheduler's code.

#### lightweight simulator

#### Assumptions

1. Lightweight simulator initializes cluster state using task-size data and only instantiates sufficiently many tasks to utilize about 60% of cluster resources.

2. We only consisder two types of jobs, *batch* and *service*

#### Metrics
1. t <sub>decision</sub> = t <sub>job</sub> + t <sub>task</sub> * *tasks per job* 

    Among the equation, t <sub>job</sub> is the overhead per job. 
2. λ <sub>jobs</sub> 
    
    It indicates job arrival rate.
3. job wait time

    difference between job submission time and the beginning of the job's first scheduling attempt.

4. scheduler busyness

5. conflict fraction

    denotes average job experiences per successful transaction