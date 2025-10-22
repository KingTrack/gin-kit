package conf

// OffsetNewest stands for the log head offset, i.e. the offset that will be
// assigned to the next message that will be produced to the partition. You
// can send this to a client's GetOffset method to get this offset, or when
// calling ConsumePartition to start consuming new messages.
// OffsetNewest int64 = -1
// OffsetOldest stands for the oldest offset available on the broker for a
// partition. You can send this to a client's GetOffset method to get this
// offset, or when calling ConsumePartition to start consuming from the
// oldest offset that is still available on the broker.
// OffsetOldest int64 = -2

type Config struct {
	Name                         string   `toml:"name" json:"name" yaml:"name"`
	Addrs                        []string `toml:"addrs" json:"addrs" yaml:"addrs"`
	GroupID                      string   `toml:"group_id" json:"group_id" yaml:"group_id"`
	OffsetsInitial               int64    `toml:"offsets_initial" json:"offsets_initial" yaml:"offsets_initial"`
	OffsetsAutoCommitEnable      bool     `toml:"offsets_auto_commit_enable" json:"offsets_auto_commit_enable" yaml:"offsets_auto_commit_enable"`
	OffsetsAutoCommitIntervalSec int      `json:"offsets_auto_commit_interval_sec" json:"offsets_auto_commit_interval_sec" yaml:"offsets_auto_commit_interval_sec"`
}
