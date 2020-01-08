# 在AWS上部署K8S集群并安装Dashboard、Jenkins

## K8S

### 部署环境

3台AWS m5.xlarge实例，一主两从配置，hostname分别为k8s-master,k8s-node1,k8s-node2。

### 部署过程

1. 关闭三台实例的swap：

    ```(sh)
    swapoff -a
    ```

2. 在三台机器上安装docker：

    ```(sh)
    apt-get update -y
    apt-get install apt-transport-https -y
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
    apt-get install docker.io -y
    # 启动docker服务
    systemctl enable docker
    systemctl start docker
    systemctl status docker
    ```

    *k8s初始化会警告：`[WARNING IsDockerSystemdCheck]: detected "cgroupfs" as the Docker cgroup driver. The recommended driver is "systemd".`, 可以通过如下解决：
    1. 修改或创建/etc/docker/daemon.json，加入下述内容：

        ```(json)
        {
            "exec-opts": ["native.cgroupdriver=systemd"]
        }
        ```

    2. 重启docker： `systemctl restart docker`

3. 在三台机器上安装k8s：

    ```(sh)
    sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    echo 'deb http://apt.kubernetes.io/ kubernetes-xenial main' | sudo tee /etc/apt/sources.list.d/kubernetes.list
    apt-get update
    apt-get install -y kubelet kubeadm kubectl kubernetes-cni
    ```

4. k8s-master节点初始化：

    ```(sh)
    #初始化，最后显示的命令是node加入集群指令
    kubeadm init --pod-network-cidr 10.244.0.0/16

    #配置kubectl
    mkdir -p $HOME/.kube
    cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    chown $(id -u):$(id -g) $HOME/.kube/config

    #安装flannel网络插件
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/k8s-manifests/kube-flannel-rbac.yml
    ```

5. node加入master:

    ```(sh)
    kubeadm join 172.31.14.182:6443 --token wzpaag.abz0nuwljizs66ti \
    --discovery-token-ca-cert-hash sha256:753e25a13676bc87597ec5b463ac51be9445cfaca606c8b170f7412713104a1c
    ```

至此，一个可用的K8S集群就搭建好了

## Dashboard

1. 部署Dashboard，采用官方文档中的配置文件：

    ```(sh)
    wget https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta4/aio/deploy/recommended.yaml
    ```

    修改service，将其中的端口暴露出来：

    ```(yaml)
    kind: Service
    apiVersion: v1
    metadata:
    labels:
        k8s-app: kubernetes-dashboard
    name: kubernetes-dashboard
    namespace: kubernetes-dashboard
    spec:
    type: NodePort
    ports:
        - port: 443
        targetPort: 8443
        nodePort: 30001
    selector:
        k8s-app: kubernetes-dashboard
    ```

2. 令牌方式登录：

    ```(sh)
    #创建serviceaccount
    kubectl create serviceaccount dashboard-serviceaccount -n kube-system

    #创建clusterrolebinding
    kubectl create clusterrolebinding dashboard-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:dashboard-serviceaccount

    #查看token列表
    kubectl get secret -n kube-system |grep dashboard-serviceaccount-token
    #输出如：
    #dashboard-serviceaccount-token-xxxxx kubernetes.io/service-account-token 3 22s

    #获取token
    kubectl describe secret dashboard-serviceaccount-token-xxxxx -n kube-system
    ```

## CI/CD

### Jenkins部署

1. 因为并非长期使用，简洁起见，没有使用RBAC，同时用emptydir挂载：

    ```(yaml)
    kind: Deployment
    metadata:
    name: jenkins
    labels:
        k8s-app: jenkins
    spec:
    replicas: 1
    selector:
        matchLabels:
        k8s-app: jenkins
    template:
        metadata:
        labels:
            k8s-app: jenkins
        spec:
        containers:
        - name: jenkins
            image: jenkins/jenkins:lts-alpine
            imagePullPolicy: IfNotPresent
            volumeMounts:
            - name: jenkins-home
            mountPath: /var/jenkins_home
            ports:
            - containerPort: 8080
            name: web
            - containerPort: 50000
            name: agent
        volumes:
            - name: jenkins-home
            emptyDir: {}

    ---

    apiVersion: v1
    metadata:
    labels:
        k8s-app: jenkins
    name: jenkins
    spec:
    type: NodePort
    ports:
        - port: 8080
        name: web
        nodePort: 30000
        - port: 50000
        name: agent
        targetPort: 50000
    selector:
        k8s-app: jenkins
    ```

    暴露了30000端口用于访问。

### Jenkins配置

1. 通过初始密码设定用户名密码和安装默认插件。
2. 安装Kubernetes插件。
3. 在Manage Jenkins中，最下面的新增一个云，添加kubernetes。
    1. kubernetes地址：`https://kubernetes.default.svc.cluster.local`
    2. Jenkins地址: `http://jenkins.kubernetes-plugin:8080`
4. 此时Jenkins已经可以使用，若需要特殊编译环境，通过添加pod template，并将container名称设置为jnlp，使得编译时拉取自定义镜像用于编译。

## 一些坑

1. AWS安全组一定要设置对，曾经发现即使配置了flannel之后网络仍然无法正确运行（表现为jenkins无法建立slave pod，测试发现dns有问题），发现是某端口（应该是3379）使用UDP协议，但是安全组配置是允许TCP流量。修改后则正确运行。

2. 一开始部署Dashboard时采用了祖传yaml文件，结果登录之后一直404，查看官方文档之后发现最新版本的k8s与旧版本的Dashboard不再兼容，使用新版之后解决问题。
