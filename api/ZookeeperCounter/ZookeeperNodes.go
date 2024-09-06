package zookeepercounter

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/SubhamMurarka/Schotky/Config"
	"github.com/go-zookeeper/zk"
)

var (
	ParentPath        = Config.Cfg.ParentPath
	GlobalCounterPath = Config.Cfg.GlobalCounterPath
	ParentLockPath    = Config.Cfg.ParentLockPath
	ChildLocks        = Config.Cfg.ChildLocks
	ServerPath        = Config.Cfg.ServerPath
	ServerPort        = Config.Cfg.ServerPort
	ServerRangePath   = Config.Cfg.ServerPath + ServerPort
	Servers           = Config.Cfg.Servers
	SessionTimeout    = Config.Cfg.SessionTimeout
	CounterRange      = Config.Cfg.CounterRange
)

// ZooKeeperClient defines the interface that will be implemented by ZKClient.
type ZooKeeperClient interface {
	Connect() error
	CreatePersistentNodes()
	GetNewRange() (int64, int64)
	Close()
}

// ZKClient is the struct that contains the ZooKeeper connection.
type ZKClient struct {
	conn *zk.Conn
	mu   sync.Mutex
}

// NewZooKeeperClient is the constructor that initializes the ZKClient struct and returns it as a ZooKeeperClient interface.
func NewZooKeeperClient() ZooKeeperClient {
	return &ZKClient{}
}

// Connect establishes a connection to the ZooKeeper server.
func (zkc *ZKClient) Connect() error {
	var err error
	zkc.conn, _, err = zk.Connect([]string{"zoo1:2181", "zoo2:2182", "zoo3:2183"}, SessionTimeout)
	if err != nil {
		log.Fatalf("Failed to connect to ZooKeeper: %v", err)
		return err
	}
	return nil
}

// CreatePersistentNodes creates persistent nodes required for the system.
func (zkc *ZKClient) CreatePersistentNodes() {
	fmt.Println(ParentPath)
	fmt.Println(GlobalCounterPath)
	fmt.Println(ParentLockPath)
	fmt.Println(ServerPath)
	fmt.Println(ServerPort)
	fmt.Println(ServerRangePath)
	createNode(zkc.conn, ParentPath, nil, 0)
	createNode(zkc.conn, GlobalCounterPath, []byte("1"), 0)
	createNode(zkc.conn, ParentLockPath, nil, 0)
	createNode(zkc.conn, ServerPath, nil, 0)
	createEphemeralNode(zkc.conn, ServerRangePath, nil)
}

// GetNewRange locks the node, updates the global counter, and returns a new range.
func (zkc *ZKClient) GetNewRange() (int64, int64) {
	// getting the mutex locke
	zkc.mu.Lock()
	defer zkc.mu.Unlock()

	// Lock the node
	path := zkc.lockNode()
	defer zkc.ReleaseNode(path)

	// Fetch and update the global counter
	data, _, err := zkc.conn.Get(GlobalCounterPath)
	if err != nil {
		log.Fatalf("Failed to get global counter: %v", err)
	}

	fmt.Println("current global counter value 1: ", data)

	currentValue, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert counter value to integer: %v", err)
	}

	fmt.Println("current global counter value 2: ", currentValue)

	CounterRange, err := strconv.ParseInt(Config.Cfg.CounterRange, 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert CounterRangePath to int64: %v", err)
	}

	fmt.Println("current counter range value: ", CounterRange)

	newValue := currentValue + CounterRange + 1
	fmt.Printf("Updating counter from %d to %d\n", currentValue, newValue)

	newData := []byte(strconv.FormatInt(newValue, 10))

	fmt.Println("new counter value is:", newData)

	// Update the global counter node
	_, err = zkc.conn.Set(GlobalCounterPath, newData, -1)
	if err != nil {
		log.Fatalf("Failed to update global counter node: %v", err)
	}

	// Update the server range node with the new range
	zkc.UpdateCounterRangeOfServers(currentValue)

	// Fetch the updated range
	data, _, err = zkc.conn.Get(ServerRangePath)
	if err != nil {
		log.Fatalf("Failed to get server range: %v", err)
	}

	rangeString := string(data)
	parts := strings.Split(rangeString, "-")
	if len(parts) != 2 {
		log.Fatalf("Invalid range format: %s", rangeString)
	}

	startValue, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert start value to integer: %v", err)
	}

	endValue, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert end value to integer: %v", err)
	}

	return startValue, endValue
}

// Close closes the ZooKeeper connection.
func (zkc *ZKClient) Close() {
	zkc.conn.Close()
	fmt.Println("ZooKeeper connection closed.")
}

// Helper function to create a node in ZooKeeper.
func createNode(conn *zk.Conn, path string, data []byte, flags int32) {
	nodePath, err := conn.Create(path, data, flags, zk.WorldACL(zk.PermAll))
	if err != nil {
		if err == zk.ErrNodeExists {
			fmt.Println("Node Already Exists:", path)
		} else {
			log.Fatalf("Failed to create node: %v", err)
		}
	} else {
		fmt.Println("Created node:", nodePath)
	}
}

// createEphemeralNode creates an ephemeral node in ZooKeeper.
func createEphemeralNode(conn *zk.Conn, path string, data []byte) {
	nodePath, err := conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		if err == zk.ErrNodeExists {
			fmt.Println("Ephemeral Node Already Exists:", path)
		} else {
			log.Fatalf("Failed to create ephemeral node: %v", err)
		}
	} else {
		fmt.Println("Created ephemeral node:", nodePath)
	}
}

// lockNode creates an ephemeral sequential node and acquires a lock by watching the previous node.
func (zkc *ZKClient) lockNode() string {
	// Create an ephemeral sequential node
	flags := int32(zk.FlagEphemeral | zk.FlagSequence)
	path, err := zkc.conn.Create(ChildLocks, nil, flags, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("Failed to create ephemeral sequential node: %v", err)
	}
	fmt.Println("Created ephemeral sequential node:", path)

	for {
		// Get all children under the lock path
		children, _, err := zkc.conn.Children(ParentLockPath)
		if err != nil {
			log.Fatalf("Failed to get children: %v", err)
		}

		// Sort the children to find the smallest sequential znode
		sort.Strings(children)

		// Find the index of the created node in the list of children
		for i, child := range children {
			fullChildPath := ParentLockPath + "/" + child
			if fullChildPath == path {
				if i == 0 {
					// Acquired the lock since this node is the smallest
					fmt.Println("Acquired the lock with Node:", path)
					return path
				}

				// Watch the previous znode in the sequence
				prevNode := ParentLockPath + "/" + children[i-1]
				_, _, watchCh, err := zkc.conn.ExistsW(prevNode)
				if err != nil {
					log.Fatalf("Failed to set watch on %s: %v", prevNode, err)
				}

				// Wait for the previous node to be deleted
				fmt.Println("Waiting for the previous node:", prevNode)
				<-watchCh
				fmt.Println("Previous node deleted, acquiring lock now:", path)
				return path
			}
		}
	}
}

// ReleaseNode deletes the znode that holds the lock.
func (zkc *ZKClient) ReleaseNode(path string) {
	err := zkc.conn.Delete(path, -1)
	if err != nil {
		log.Fatalf("Failed to delete node %s: %v", path, err)
	}
	fmt.Println("Node released and deleted:", path)
}

// UpdateCounterRangeOfServers updates the server's range in the ZooKeeper node.
func (zkc *ZKClient) UpdateCounterRangeOfServers(currentValue int64) {
	CounterRange, err := strconv.ParseInt(Config.Cfg.CounterRange, 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert CounterRangePath to int64: %v", err)
	}

	// Format the range as a string
	rangeData := fmt.Sprintf("%d-%d", currentValue, currentValue+CounterRange)

	// Update the server range path node
	_, err = zkc.conn.Set(ServerRangePath, []byte(rangeData), -1)
	if err != nil {
		log.Fatalf("Failed to update server range: %v", err)
	}

	// Fetch the updated server range value from ZooKeeper
	data, _, err := zkc.conn.Get(ServerRangePath)
	if err != nil {
		log.Fatalf("Failed to fetch server range: %v", err)
	}

	// Print the fetched value
	fmt.Printf("Fetched server range: %s\n", string(data))

	fmt.Println("Server range updated successfully.")
}
