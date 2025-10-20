package instance

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
)

const (
	MetaKeyWeight = "weight"
)

type Instance struct {
	ServiceName string            `json:"service_name"` // 服务名字唯一标识
	Schema      string            `json:"schema"`
	IP          string            `json:"ip"`
	Port        int               `json:"port"`
	Weight      int               `json:"weight"`
	Meta        map[string]string `json:"tags"`
}

func (i *Instance) GetHost() string {
	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

func (i *Instance) IsEqualEndpoint(v *Instance) bool {
	if v == nil {
		return false
	}
	return i.IP == v.IP && i.Port == v.Port
}

func (i *Instance) GetMeta() map[string]string {
	i.Meta[MetaKeyWeight] = strconv.Itoa(i.Weight)
	return i.Meta
}

func (i *Instance) ServiceID() string {
	builder := &strings.Builder{}
	builder.WriteString(i.ServiceName)
	builder.WriteString(i.IP)
	builder.WriteString(strconv.Itoa(i.Port))
	return fmt.Sprintf("%x", md5.Sum([]byte(builder.String())))
}

func MatchMeta(condition map[string]string, input map[string]string) bool {
	for key, value := range condition {
		if x, ok := input[key]; !ok || x != value {
			return false
		}
	}
	return true
}

func GetWeight(meta map[string]string) int {
	weight, _ := strconv.Atoi(meta[MetaKeyWeight])
	return weight
}

func RebuildMeta(meta map[string]string) map[string]string {
	m := make(map[string]string, len(meta))
	for k, v := range meta {
		if k == MetaKeyWeight {
			continue
		}
		m[k] = v
	}
	return m
}
