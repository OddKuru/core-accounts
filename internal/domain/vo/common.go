package vo

type StringValuer interface {
	Value() string
	Validate() error
}

var (
	_ StringValuer = (*ID)(nil)
	_ StringValuer = (*LoginName)(nil)
	_ StringValuer = (*Email)(nil)
	_ StringValuer = (*Password)(nil)
)
