package main

import (
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
	PodFitResourcePred = "PodFitResourcePred"
)

var predicatesFuncs = map[string]FitPredicate{
	PodNameFitPred: PodNameFitPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{
	PodNameFitPred,
	PodFitResourcePred,
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
 * check if pod name length is within [max(node name length - 5, 0), min(node name length + 5, 64)]
 */
func PodNameFitPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var valid bool
	min := math.Min(float64(len(node.Name)), 32)
	max := math.Max(0, float64(len(node.Name)))
	valid = int(min) < len(pod.Name)  && len(pod.Name) < int(max)
	if valid {
		log.Printf("pod %v/%v length is %d, fit on node %v\n", pod.Name, pod.Namespace, len(pod.Name), node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v length is %d, not fit on node %v\n", pod.Name, pod.Namespace, len(pod.Name), node.Name)
	return false, []string{PodNameFitPredFailMsg}, nil
}
