package freestuff

import "strings"

type Label struct {
	Key   string
	Value string
}

type ExtraInfo interface {
	GetLabels() []Label
}

func GetExtraInfo(link string) (ExtraInfo, error) {
	if strings.Contains(link, "apple.com") {
		return GetAppstoreInfo(link)
	}

	return nil, nil
}
