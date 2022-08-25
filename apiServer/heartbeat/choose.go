package heartbeat

import "math/rand"

//被选泽的数据节点 (exclude) 不会再被选择;修复数据的时候也会用到
func ChooseRandomDataServer(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return
}
