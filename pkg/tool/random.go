package tool

import (
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	clash     ClashYaml
	poolMutex sync.RWMutex
	proxyPool []map[string]interface{}
	loadOnce  sync.Once
)

type ClashYaml struct {
	Proxy []map[string]interface{} `json:"proxies" yaml:"proxies"`
}

func RandomProxySimple(proxyFilePath string) (proxy map[string]interface{}) {
	// Ensure the proxy pool is loaded once
	loadOnce.Do(func() {
		loadProxies(proxyFilePath)
	})

	// 获取读锁
	poolMutex.RLock()
	if len(proxyPool) == 0 {
		// 释放读锁
		poolMutex.RUnlock()
		// 获取写锁
		poolMutex.Lock()
		if len(proxyPool) == 0 {
			// Refill the pool and shuffle
			proxyPool = append(proxyPool, clash.Proxy...)
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(proxyPool), func(i, j int) { proxyPool[i], proxyPool[j] = proxyPool[j], proxyPool[i] })
		}
		// 释放写锁
		poolMutex.Unlock()
		// 获取读锁
		poolMutex.RLock() // Downgrade back to read lock
	}

	// Get a proxy from the pool
	proxy, proxyPool = proxyPool[len(proxyPool)-1], proxyPool[:len(proxyPool)-1]
	// 释放读写锁
	poolMutex.RUnlock()
	return proxy
}

func loadProxies(proxyFilePath string) {
	filepath.Walk(proxyFilePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			content, _ := os.ReadFile(path)
			p := ClashYaml{}
			yaml.Unmarshal(content, &p)
			clash.Proxy = append(clash.Proxy, p.Proxy...)
		}
		return nil
	})
	proxyPool = append(proxyPool, clash.Proxy...)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxyPool), func(i, j int) { proxyPool[i], proxyPool[j] = proxyPool[j], proxyPool[i] })
}

func RandomProxy(cfg string) (str map[string]interface{}, cur int, total int) {
	clash := ClashYaml{}
	filepath.Walk(cfg, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			content, _ := os.ReadFile(path)
			p := ClashYaml{}
			yaml.Unmarshal(content, &p)
			clash.Proxy = append(clash.Proxy, p.Proxy...)
		}
		return nil
	})
	l := len(clash.Proxy)
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(l)
	if randomInt >= l {
		randomInt = l - 1
	}
	return clash.Proxy[randomInt], randomInt, l
}
