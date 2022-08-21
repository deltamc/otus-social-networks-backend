package tar

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"os"
	"sync"
	"time"
)

var tarantoolConnect *tarantool.Connection

func Check(n string) bool {
	client := Client(n)
	if client == nil {
		return false
	}
	_, err := client.Ping()
	if err != nil {
		fmt.Println("tarantool", err.Error())
		return false
	}
	return true
}

var lockTarantool sync.RWMutex

func Client(n string) *tarantool.Connection {
	lockTarantool.Lock()
	defer lockTarantool.Unlock()

	if tarantoolConnect != nil {
		return tarantoolConnect
	}

	var err error
	fmt.Println(os.Getenv("TARANTOOL_" + n + "_USER"))
	opts := tarantool.Opts{
		Reconnect:     1 * time.Second,
		MaxReconnects: 5,
		User:          os.Getenv("TARANTOOL_" + n + "_USER"),
		Pass:          os.Getenv("TARANTOOL_" + n + "_PASSWORD"),
	}

	tarantoolConnect, err = tarantool.Connect(fmt.Sprintf("%s:%s",
		os.Getenv("TARANTOOL_"+n+"_HOST"),
		os.Getenv("TARANTOOL_"+n+"_PORT"),
	), opts)

	if err != nil {
		fmt.Println("Connection refused:", err)
	}

	return tarantoolConnect
}

func Close() {
	if tarantoolConnect == nil {
		return
	}
	err := tarantoolConnect.Close()
	if err != nil {
		fmt.Println("Connection close:", err)
	}
}
