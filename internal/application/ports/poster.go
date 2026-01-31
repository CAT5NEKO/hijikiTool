package ports

type Poster interface {
	Post(content string) error
}
