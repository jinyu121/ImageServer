package filter

const (
	PREDEFINE_DEFAULT   = "@DEFAULT"
	PREDEFINE_IMAGE_EXT = "@IMAGE"
	PREDEFINE_VIDEO_EXT = "@VIDEO"
	PREDEFINE_AUDIO_EXT = "@AUDIO"
)

type Filter interface {
	Filter(string) bool
	Extract(string) []string
}
