package topic

type Sort string

const (
	SortActive  Sort = "active"
	SortNew     Sort = "new"
	SortReplies Sort = "replies"
)

func Sorts() []string {
	return []string{string(SortActive), string(SortNew), string(SortReplies)}
}
