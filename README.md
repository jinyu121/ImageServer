# Image Server

Image Server sets up a simple image server, to let you view images.

## How to use

Open your browser and navigate to `http://your_ip:9420` to view your images.

### Browse the current folder

```shell
ImageServer
```

### Browse a specified folder, list, or LMDB file

```shell
ImageServer path/to/your/folder
ImageServer path/to/your/list.txt
ImageServer path/to/your/list.csv --column 0
ImageServer path/to/your/lmdb/database.lmdb
```

### Limit page size

```shell
ImageServer path/to/your/folder --page 100
```