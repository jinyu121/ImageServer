# Image Server

Image Server sets up a simple HTTP service, to make you view batch of images easier at a glance.

Start the server, open your browser, and navigate to `http://your_ip:9420` to view your images.

## How to start the server

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

## How to browse in the browser

### Basic usage

```
http://your_ip:9420
```

### Pagination

```
http://your_ip:9420?p=2333
```

### Compare folders

```
http://your_ip:9420/path/to/foder/one?c=path/to/folder/two
http://your_ip:9420/path/to/foder/one?c=path/to/folder/two&c=/path/to/another/folder
```

## UI

### Navigation

- ⏮ : Previous Neighborhood Folder
- ⏫ : Parent Folder
- ⏭ : Next Neighborhood Folder

### Pagination

- 🔼 : Previous Page
- 🔽 : Next Page
