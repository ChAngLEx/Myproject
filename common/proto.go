package common

//StdNew 标准格式
type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}

//NewsSource URL:API Parse:API结果转StdNews的函数
type NewsSource struct {
	URL   string
	Parse func([]byte) ([]StdNew, error) //函数可以作为变量
}
