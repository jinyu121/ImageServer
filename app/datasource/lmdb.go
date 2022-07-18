package datasource

import (
	"log"
	"path"
	"sort"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/schollz/progressbar/v3"
)

type LmdbDataSource struct {
	Root string

	lmdbEnv      *lmdb.Env
	lmdbDbi      lmdb.DBI
	lmdbFileTree *Node
}

func NewLmdbDataSource(filePath string) *LmdbDataSource {
	ds := &LmdbDataSource{Root: "/"}
	ds.scan(filePath)
	return ds
}

func (ds *LmdbDataSource) scan(filePath string) {
	ds.lmdbFileTree = NewRoot()

	ds.lmdbEnv, _ = lmdb.NewEnv()
	_ = ds.lmdbEnv.SetMaxDBs(1)
	_ = ds.lmdbEnv.SetMapSize(1 << 42)
	err := ds.lmdbEnv.Open(filePath, 0, 644)

	if nil != err {
		panic(err)
	}
	err = ds.lmdbEnv.Update(func(txn *lmdb.Txn) (err error) {
		ds.lmdbDbi, err = txn.OpenRoot(lmdb.Create)
		return
	})
	if nil != err {
		panic(err)
	}

	log.Printf("Scan database %s", filePath)

	bar := progressbar.Default(-1, "Scanning")
	callback := func(itemPath string) {
		_ = bar.Add(1)
	}

	counter := 0
	_ = ds.lmdbEnv.View(func(txn *lmdb.Txn) (err error) {
		cur, err := txn.OpenCursor(ds.lmdbDbi)
		if err != nil {
			return err
		}
		defer cur.Close()

		for {
			k, _, err := cur.Get(nil, nil, lmdb.Next)
			if lmdb.IsNotFound(err) {
				return nil
			}
			if err != nil {
				return err
			}

			ds.lmdbFileTree.Add(string(k))
			callback(string(k))
			counter += 1
		}
	})
	log.Printf("Scan Done! %d records scan", counter)
}

func (ds *LmdbDataSource) GetFile(filePath string) (data []byte, err error) {
	_ = ds.lmdbEnv.View(func(txn *lmdb.Txn) (err error) {
		data, err = txn.Get(ds.lmdbDbi, []byte(filePath[1:]))
		return nil
	})

	return data, err
}

func (ds *LmdbDataSource) GetFolder(current string) (content FolderContent, err error) {
	content = FolderContent{
		Name:    current,
		Folders: []string{},
		Files:   []string{},
	}

	node, err := ds.lmdbFileTree.GetChild(current)
	if nil != err {
		return content, err
	}

	basePath := node.GetAbsolutePath()
	for _, v := range node.Children {
		pa := path.Join(basePath, v.Name)
		if v.IsFile {
			content.Files = append(content.Files, pa)
		} else {
			content.Folders = append(content.Folders, pa)
		}
	}

	// Sort to keep a static order
	sort.Strings(content.Folders)
	sort.Strings(content.Files)

	return content, nil
}

func (ds *LmdbDataSource) GetNeighbor(current string) (nav *Navigation) {
	nav = &Navigation{}
	if "/" == current || "" == current {
		return
	}

	node, err := ds.lmdbFileTree.GetChild(current)
	if nil != err {
		return
	}
	nav.Current = node.GetAbsolutePath()

	parent := node.Parent
	if nil != node.Parent {
		nav.Parent = parent.GetAbsolutePath()
	}

	folders := make([]*Node, 0)
	for _, v := range parent.Children {
		if !v.IsFile {
			folders = append(folders, v)
		}
	}
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Name < folders[j].Name
	})

	for i, val := range folders {
		if val.Name == node.Name {
			if i-1 >= 0 {
				nav.Prev = folders[i-1].GetAbsolutePath()
			}
			if i+1 < len(folders) {
				nav.Next = folders[i+1].GetAbsolutePath()
			}
			return
		}
	}
	return
}

func (ds *LmdbDataSource) Stat(filePath string) *FileStat {
	result := &FileStat{
		Exists: false,
		IsFile: false,
	}

	if node, err := ds.lmdbFileTree.GetChild(filePath); nil == err {
		result.Exists = true
		result.IsFile = node.IsFile
	}

	return result
}
