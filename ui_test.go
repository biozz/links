package main

import (
	"io/fs"
	"regexp"
	"testing"

	"github.com/pocketbase/pocketbase/ui"
)

// Check the embedded filesystem used by the Superuser UI, not a local source directory.
func TestSuperuserUIAssets(t *testing.T) {
	index, err := fs.ReadFile(ui.DistDirFS, "index.html")
	if err != nil {
		t.Fatal(err)
	}
	refs := regexp.MustCompile(`(?:src|href)="\./([^"]+)"`).FindAllSubmatch(index, -1)
	if len(refs) == 0 {
		t.Fatal("Superuser UI has no asset references")
	}
	for _, ref := range refs {
		path := string(ref[1])
		t.Run(path, func(t *testing.T) {
			info, err := fs.Stat(ui.DistDirFS, path)
			if err != nil {
				t.Fatal(err)
			}
			if info.IsDir() || info.Size() == 0 {
				t.Fatalf("invalid embedded asset: %s", path)
			}
		})
	}
}
