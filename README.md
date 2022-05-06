# Image Server

Image Server sets up a simple HTTP service, to make you view batch of images easier.

## How to use

Start the server, open your browser and navigate to `http://your_ip:9420` to view your images.

### Browse the current folder

```shell
ImageServer
```

### Browse a specified folder, list, or LMDB file

```shell
ImageServer path/to/your/folder
ImageServer path/to/your/list.txt
ImageServer --column 0 path/to/your/list.csv 
ImageServer path/to/your/lmdb/database.lmdb
```

### Limit page size

```shell
ImageServer --page 100 path/to/your/folder
```

### Change port

Default port is 9420, you can change it by `--port` option.

```shell
ImageServer --port 2333 path/to/your/folder
```