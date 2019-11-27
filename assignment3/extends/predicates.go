package main

import (
	"fmt"
	_ "fmt"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"log"
	"math"
	"strings"
)

const (
	// LuckyPred rejects a node if you're not lucky ¯\_(ツ)_/¯
	PodNameFitPred        = "PodNameFitPred"
	PodNameFitPredFailMsg = "Pod name don't fit with node"
)

var predicatesFuncs = map[string]FitPredicate{
	PodNameFitPred: PodNameFitPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{
	PodNameFitPred,
}

// filter filters nodes according to predicates defined in this extender
// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFit()
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod

	// TODO: parallelize this
	// TODO: hanlde error
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	failReasons := []string{}
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

/**
 * check if pod name length is within [0, max(node.index + 17, 32))
 * node_index is 1, 2, 3
 */
func PodNameFitPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var valid bool
	index := int(node.Name[9] - 0x30)
	max := math.Min(float64(index) + 17, 32)
	valid = int(max) > len(pod.Name)
	if valid {
		log.Printf("[%v] pod %v/%v length is %d, node length is %d, fit", node.Name, pod.Name, pod.Namespace, len(pod.Name), int(max))
		return true, nil, nil
	}
	log.Printf("[%v] pod %v/%v length is %d,  node length is %d, unfit\n", node.Name, pod.Name, pod.Namespace, len(pod.Name), int(max))
	return false, []string{PodNameFitPredFailMsg}, fmt.Errorf("pod length exceed ")
}
