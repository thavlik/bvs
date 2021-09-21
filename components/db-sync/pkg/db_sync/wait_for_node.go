package db_sync

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/thavlik/bvs/components/node/pkg/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForNode(cl client.Client) error {
	if err := WaitForNodeK8s(cl); err != nil {
		return fmt.Errorf("WaitForNodeK8s: %v", err)
	}
	if err := WaitForNodeHTTP(0); err != nil {
		return fmt.Errorf("WaitForNodeHTTP: %v", err)
	}
	return nil
}

func RequireEnv(name string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("missing environment variable '%s'", name))
	}
	return v
}

func WaitForNodeK8s(cl client.Client) error {
	name := RequireEnv("POD_NAME")
	namespace := RequireEnv("POD_NAMESPACE")
	timeout := 8 * time.Minute
	start := time.Now()
	stop := make(chan int, 1)
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(30 * time.Second):
				fmt.Printf("Still waiting for node and postgres containers to be Ready (%s elapsed)\n",
					time.Since(start).Round(time.Second).String())
			}
		}
	}()
	defer func() {
		stop <- 1
		close(stop)
	}()
	fmt.Println("Waiting for node and postgres containers to be Ready")
	for time.Since(start) < timeout {
		pod := &v1.Pod{}
		if err := cl.Get(context.TODO(), types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}, pod); err != nil {
			return fmt.Errorf("k8s: %v", err)
		}
		nodeReady := false
		postgresReady := false
		for _, container := range pod.Status.ContainerStatuses {
			if container.Ready {
				if container.Name == "node" {
					nodeReady = true
					continue
				}
				if container.Name == "postgres" {
					postgresReady = true
					continue
				}
			}
		}
		if nodeReady && postgresReady {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("timed out after %v", timeout)
}

func WaitForNodeHTTP(timeout time.Duration) error {
	cl := api.NewNodeClient("http://localhost:80", "", "", 10*time.Second)
	if timeout == 0 {
		resp, err := cl.ProbeReady(context.TODO(), api.ProbeReadyRequest{})
		if err != nil {
			return fmt.Errorf("node: %v", err)
		}
		if resp.IsReady {
			return nil
		}
		return fmt.Errorf("node is not synchronized and timeout is zero")
	}
	start := time.Now()
	stop := make(chan int, 1)
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(30 * time.Second):
				fmt.Printf("Still waiting for node to sync with network\n")
			}
		}
	}()
	defer func() {
		stop <- 1
		close(stop)
	}()
	fmt.Println("Ensuring node is fully synchronized")
	for time.Since(start) < timeout {
		resp, err := cl.ProbeReady(context.TODO(), api.ProbeReadyRequest{})
		if err != nil {
			return fmt.Errorf("node: %v", err)
		}
		if resp.IsReady {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("timed out after %v", timeout)
}
