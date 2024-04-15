package supply

import (
	"io"
	"os"
	"sf-mu/pkg/models/censors"

	"strings"
)

func BlackList() ([]censors.BlackList, error) {
	f, err := os.Open("./cmd/censors/pkg/supply/words.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	var sl []censors.BlackList
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		str := censors.BlackList{
			BanWord: trimmedLine,
		}
		sl = append(sl, str)
	}

	return sl, nil
}
