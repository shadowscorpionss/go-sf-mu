package censors

type BlackList struct {
	ID      int    `json:"ID,omitempty"`
	BanWord string `json:"banWord,omitempty"`
}
//censors contract
type Interface interface {
	AllList() ([]BlackList, error)
	AddList(c BlackList) error
	CreateBlackListTable() error
	DropBlackListTable() error
}
