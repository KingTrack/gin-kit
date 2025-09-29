package define

import (
	"github.com/KingTrack/gin-kit/kit/types/datacenter/instance"
)

type IBalancer interface {
	Pick(meta map[string]string, skip *instance.Instance) (*instance.Instance, error)
	Update(instances []instance.Instance)
}
