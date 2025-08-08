package mapper

const (
	DefaultLangCode = "id-ID"
	TimeFormat      = "2006-01-02 15:04:05"
	DateFormat      = "2006-01-02"
)

func Ptr[T any](v T) *T {
	return &v
}
