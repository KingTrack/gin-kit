package conf

import (
	"fmt"
	"os"

	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/plugin/source"
)

type Resource struct{}

func NewSource() source.ISource {
	return &Resource{}
}

func (r *Resource) Load(path string) (map[string][]byte, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	return map[string][]byte{tlscontext.Default.ToString(): data}, nil
}

func (r *Resource) Watch(rootPath string, callback func(namespaceKey string, data []byte)) {
	return
}
