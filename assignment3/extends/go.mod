module extends

go 1.12

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191115221445-ec04ad4dbd24

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20191115220130-ea09a2678486

replace k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191114223734-b81ee8851698

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191116142037-570483aa5fd4

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191116203941-08e4eafd6d11

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20191115221332-dea8d0e6b550

replace k8s.io/api => k8s.io/api v0.0.0-20191115135540-bbc9463b57e5

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191115222459-ca002f3fd118

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191114215150-2a85f169f05f

replace k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191115220603-3ef46478e636

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191114224219-ae7460d41375

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191115221917-22ad31b04cf4

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191115222239-d1c34eac545e

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191117020858-b615a37f53e7

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191115222025-96c7256ae8d3

replace k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191115222826-fbc5d36fee2d

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191115220737-d15885da1fea

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191117060219-524f257d2214

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191115221700-db89989543a2

replace k8s.io/cri-api => k8s.io/cri-api v0.0.0-20191114225051-381f557272b4

require (
	github.com/julienschmidt/httprouter v1.3.0
	k8s.io/api v0.0.0
	k8s.io/kubernetes v1.16.3
)
