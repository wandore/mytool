example

consistent hash:
```
/*
learn from stathat.com/c/consistent
*/
package main

import (
    "fmt"
    uuid "github.com/satori/go.uuid"
    "github.com/wandore/mytool/hash"
)

var replicas = 1024
var clusterNum = 10
var nodes = make([]string, 0)
var nums = make([]int, 0)

func main() {
    hash := hash.New(replicas, nil)

    for i := 0; i < clusterNum; i++ {
    	node := uuid.NewV4().String()
    	hash.Add(node)
    	nodes = append(nodes, node)
    	nums = append(nums, 0)
    }

    for i := 0; i < 1000; i++ {
    	key := uuid.NewV4().String()
    	node, _ := hash.Match(key)
    	for j, member := range nodes {
    	    if node == member {
                nums[j]++
                break
            }
        }
    }

    fmt.Println(nums)

    removeNode := nodes[clusterNum-1]
    nodes = nodes[:clusterNum]
    fmt.Println("will remove", removeNode)
    hash.Remove(removeNode)
    nums = nums[:clusterNum-1]
    for i := range nums {
    	nums[i] = 0
    }

    for i := 0; i < 1000; i++ {
    	key := uuid.NewV4().String()
    	node, _ := hash.Match(key)
    	for j, member := range nodes {
    		if node == member {
    			nums[j]++
    			break
    		}
    	}
    }

    fmt.Println(nums)
}
```