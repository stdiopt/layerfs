# LayerFS

Go fs.FS layering and mounting.

## Example

```go
package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/stdiopt/layerfs"
)

//go:embed assets
var assets embed.FS

func main() {
	mfs := layerfs.FS{
		assets,
		layerfs.Prefix("/assets", os.DirFS("./data")),
	}
	sfs, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatal(err)
	}
	mfs.Mount("other", sfs)

	err = fs.WalkDir(mfs, ".", func(p string, _ fs.DirEntry, _ error) error {
		fmt.Println(p)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	f, err := mfs.Open("/assets/myfile2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err = io.Copy(os.Stdout, f); err != nil {
		log.Fatal(err)
	}

	// result:
	//
	// .
	// other
	// other/myfile1.txt
	// other/myfile2.txt
	// assets
	// assets/myfile2.txt
	// assets/myfile3.txt
	// assets/myfile1.txt

	// overriden file2 content.
}
```
