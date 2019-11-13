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

### lightweight simulator

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

scheduler busyness and conflict fraction are medians of the daily values,
while wait time values are overall averages.

#### Scheduler Performance

To better demonstrate Omega's advantage, this paper compare three types of schedulers with the above metrics
First, we will show monolithic schedulers' perfromance, then two-level schedulers', finally Omega's.

In the following figures, A、B、C represents three clusters. solid lines represent batch jobs, while 
dashed lines represent service jobs.

##### Monolithic scheduler

Monolithic scheduler is a baseline compared with the other two types of schedulers.
Monolithic schedulers are compares with the same decision time for jobs, from which we can see the need for 
partitioning job types.</br>
In this paper, authors also privide a multi-path monolithic schedulers, which equipped with a fast code
path for batch jobs</br>

X-axis is increased by t <sub>job</sub>, time for job decision.

![](./imgs/scheduler-cmp-wait-1.png)


![](./imgs/scheduler-cmp-busyness-1.png)

Results show that multi-path scheduler's performance is better thant single-path
scheduler in this situation. That's exactly true in other cases, since multi-path 
develop code path for specialized jobs, while single-path scheduler server each job
as the same. Besides, we can see that scheduler's mean job wait time is positively
related to scheduler busyness.

##### Two-level scheduler

In this paper, authors simplified test by making some assumptions.

1. scheduler only looks at the set of resources available to it 
when it begins a scheduling attempt for a job. (Actually, it is the assumption
of two-level scheduler in this paper.)

2. resource manager algorithm is really fast, so we will assume it takes 1 ms.

3. We choose a typical two-level scheduler, Mesos. Use two schedulers, batch scheduler 
framework and service scheduler framework

The results are shown below.

![](./imgs/scheduler-cmp-twolevel-1.png)

From the results, we can see by adjusting t<sub>job</sub>, the scheduler busyness turns out
to be much higher than in the monolithic multi-path case. 

As the authors explianed, it is because the concurrency control using in Mesos offer model.
Mesos uses a pessimistic concurrency control. Mesos pratitions workloads. In this example,
Mesos mainly deals with batch and service jobs. It assume that batch jobs are fast decided 
and resources are frequently available. When all of its assumptions are correct(I think in
normal cases it is correct), the pessimistic concurrency control which locks nearly lock all cluster
resources is not expensive. However, in this case, due to long decision time, the cost for waiting
is rather expensive.

We can also see another measure on abandoned jobs. This is because long-runnig service jobs took up too
much resources so that there was no sufficient resource for one batch job. While service jobs 
was running, batch scheduler try it best to get allocated. Once failed, it will not stop but retry. Unfortunately,
long-running service jobs would not release the resource, as a result, after 1000 times 
retry limit, batch scheduler would abandon the job.

##### Omega

As the same with Two-level scheduler, for Omega, we also simulate two schedulers, batch job 
scheduler and service job scheduler. We use incremental conflict fraction, which means only those changes
that do not affect an committed machine will be accepted in one transaction.

![job wait time](./imgs/scheduler-cmp-shared-wait-1.png)

![scheduler busyness](./imgs/scheduler-cmp-busyness-1.png)

The above figure shows job wait time and busyness of Omega. We can see that it is comparable with multi-path monolithic
schedulers, which means conflicts and interference are relatively rare.

In this workload, Omega works just like multi-path monolithic schedulers but do not suffer from
head of line blocking for the complete parallel scheduling. Unlike two-level schedulers, conflicts are
rare and job abandon is rare. Omega schedulers hold entire cluster state and can make good decisions.

The authors also simulate the shcedulers' scalality. In ths experiment, λ <sub>jobs</sub> represent
batch jobs arrival rate. 

In these figures, dashed verticl line indicates that the corresponding cluster scale to that times
of original workload.

![vary job rate](./imgs/omega-lightweight-jobrate-1.png)

From the results, we can see that batch scheduler is more busy then before and due to 
batch scheduler busyness, both batch jobs wait time and service job wait time increased.
According to the authors, batch jobs' wait time increasement are mainly because high job arrival rate,
while service jobs' account to additional conflicts.

batch scheduler is the main bottleneck of this workload(it is obvious since we only change the batch job
rate. this is also compatible to real world cases.). The authors also did a research on multiple batch 
scheduler. They add multiple batch job schedulers into this experiment.

![scale up batch scheduler](./imgs/omega-lightweight-batchscale.png)

In multiple batch schedulers, the authors used a simple hashing function to determine 
which scheduler to use. Yes, the multiple batch schedulers are load-balanced.

The results shows that the more batch schedulers, the more conflict factions.
However, scheduler busyness figure relects that batch schedulers are not as busy as before. It means
that conflict fractions are compensated by more schedulers.

### high-fidelity simulator

In ths experiment, the authors take research on service shceduler busyness varies as a function of 
both t <sub>job</sub> and t <sub>task</sub>.</br> 

![service job busyness](./imgs/omega-service-busyness-1.png)

For the first figure, it is a month-long trace of cluster C.</br>
We can see despite the increasing of t <sub>job</sub> and t <sub>task</sub> 
for a month-long trace of cluster C,service scheduler busyness remains low. It is obvious I think,
since long-running service jobs take up more resource and decision time. Despite the time to
decision is increasing, the relatice relative increase is not so high.

![omega service scheduler detail](./imgs/omega-service-2.png)

For the second figure, we dive into service busyness in detail.</br>
The results show that as service jobs decision time increasing, conflic fraction increases.
At about 10s, conflict fraction reaches 1.0 which means one successfule allocation requires
at least one retry.</br>
For figure c of the second figure, it indicates that conflic fraction makes it
almost 40% percent more busy that no conflict service scheulers.

The authors also use the high-fiedlity simulator to explore two implementations choices they 
were considering for Omega.
1. coarse -grained conflict detection
    Reject placement choice if any changes had been made to the targe tmachine since the local 
    copy of cell state was synchronized at the beginning of transaction.
2. all-or-nothing scheduing
    As explianed before, an entire cell state transaction would be rejected if iw would cause 
    any machine to be over-committed.

Just like the method names indicate, coarse-grained conflict detection makes things worse.
Al-or-nothing should be used only on a job-level granularity.

## Specialized scheduler, MapReduce job

### What is MapReduce

MapReduce is a method to map a job into separated jobs and then reduce the jobs into one result.
It is typically used in distributed systems and assign jobs to each system.

### Application of Omega 

Basic idea is explained in this paper.</paper>

The MapReduce scheduler(a scheduler specilized for MapReduce job) automatically choose
idle cluster resources to speed up MapReduce jobs.

> It observes the overall resource utilization in the cluster, predicts the benefits of
scaling up current and pending MapReduce jobs, and apportions some fraction of the unused resources across those jobs
according to some policy.

In this paper, MapReduce Scheduler's evaluation is simple.</br>
They also talked about implementation. Since MapReduce jobs typically have many more of activities than configured 
workers(end of distributed system may down, thus need reassigned to another end computer), they usually run out of 
available resources. Thus, the authors referred to three methods
1. max-parallelism
    
    which keeps on adding workers as long as benefit is obtained.
2. global cap
    
    which stops the MapReduce scheduler using idle resources if the total cluster utilization is above a target value
3. relative job size
    
    which limits the maximum number of workers to four times as many as it initially requested.

### Evaluation

![Three implementations of MapReduce](./imgs/omega-mapreduce-2.png)

The first figure shows the advantage using Omega scheduler. We can see as 
we add more additional resources, speedup is higher. Max-parallel policy's 
performance is typically better than the other two policies.

## My own thinking

1. 要努力发现技术和现在实际情况之间的矛盾，为什么越来越慢，一定要找到每个技术的本质上的问题，才能对症下药。
2. 通常，采用将几个技术的优势结合起来的新技术都会效果显著。
3. 技术不分类别，在数据库中使用额optimistic concurrency control在调度中也能用，在File System中更加能用，
对于新的技术多加关注，联想其他应用场景。
4. 不同的workload适用不同的技术和方法，没有能够用一种方法适应全部workload的。

