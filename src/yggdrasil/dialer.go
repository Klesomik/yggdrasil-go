package yggdrasil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/yggdrasil-network/yggdrasil-go/src/crypto"
)

// Dialer represents an Yggdrasil connection dialer.
type Dialer struct {
	core *Core
}

// Dial opens a session to the given node. The first paramter should be "nodeid"
// and the second parameter should contain a hexadecimal representation of the
// target node ID.
func (d *Dialer) Dial(network, address string) (Conn, error) {
	conn := Conn{
		mutex: &sync.RWMutex{},
	}
	nodeID := crypto.NodeID{}
	nodeMask := crypto.NodeID{}
	// Process
	switch network {
	case "nodeid":
		// A node ID was provided - we don't need to do anything special with it
		if tokens := strings.Split(address, "/"); len(tokens) == 2 {
			len, err := strconv.Atoi(tokens[1])
			if err != nil {
				return Conn{}, err
			}
			dest, err := hex.DecodeString(tokens[0])
			if err != nil {
				return Conn{}, err
			}
			copy(nodeID[:], dest)
			for idx := 0; idx < len; idx++ {
				nodeMask[idx/8] |= 0x80 >> byte(idx%8)
			}
			fmt.Println(nodeID)
			fmt.Println(nodeMask)
		} else {
			dest, err := hex.DecodeString(tokens[0])
			if err != nil {
				return Conn{}, err
			}
			copy(nodeID[:], dest)
			for i := range nodeMask {
				nodeMask[i] = 0xFF
			}
		}
	default:
		// An unexpected address type was given, so give up
		return Conn{}, errors.New("unexpected address type")
	}
	conn.core = d.core
	conn.nodeID = &nodeID
	conn.nodeMask = &nodeMask
	conn.core.router.doAdmin(func() {
		conn.startSearch()
	})
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn, nil
}
