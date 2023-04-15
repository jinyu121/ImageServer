package filter

const (
	PredefineDefault  = "@DEFAULT"
	PredefineImageExt = "@IMAGE"
	PredefineVideoExt = "@VIDEO"
	PredefineAudioExt = "@AUDIO"
)

type Filter interface {
	Filter(string) bool
	Extract(string) []string
}
