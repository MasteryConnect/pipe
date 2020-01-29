package fs_test

import (
	"github.com/masteryconnect/pipe/extras/fs"
	"github.com/masteryconnect/pipe/line"
	"github.com/masteryconnect/pipe/message"
	"github.com/masteryconnect/pipe/x"
)

func ExampleList_P() {
	line.New().SetP(
		fs.List{
			Root:         "/",
			IncludeDirs:  true,
			ExcludeFiles: true,
			Recursive:    false,
		}.P,
	).Add(
		line.I(func(m interface{}) (interface{}, error) {
			msg := m.(message.FileInfo)
			if msg.IsDir() {
				return "**dir**", nil
			}
			return m, nil
		}),
		x.Head{N: 1}.T, // limit the list to just 1 for output testing
		line.Stdout,
	).Run()

	// Output: **dir**
}
