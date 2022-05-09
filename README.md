# Image Server

Image Server sets up a simple HTTP service, to make you view batch of images easier at a glance.

Start the server, open your browser, and navigate to `http://your_ip:9420` to view your images.

## Start the server

### Current folder

```shell
ImageServer
```

### Specified folder

```shell
image_server path/to/your/folder
```

### URL list

Each line should be a URL

```shell
image_server path/to/your/list.txt
```

### CSV / TSV file

You should specify which column (by index, not name) is the image URL

```shell
image_server --column 0 path/to/your/list.csv
```

### JSON file

One line per JSON object (ImageServer do not support human-friendly formatted JSON file), and you should specify how to
get the image URL by JSONPath syntax

```shell
image_server --json "@.images[*]" path/to/your/json/file.json  
```

### LMDB file

Keys will be separated by `/` to get virtual paths/folders

```shell
image_server path/to/your/lmdb/database.lmdb
```

### Parameters

#### Change page size

Default page size is 1000, you can change it by `--page` option

```shell
image_server --page 100 path/to/your/folder
```

#### Change port

Default port is 9420, you can change it by `--port` option

```shell
image_server --port 2333 path/to/your/folder
```

## Browse images in your browser

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

## Web UI

### Navigation

- ⏮ : Previous Neighborhood Folder
- ⏫ : Parent Folder
- ⏭ : Next Neighborhood Folder

### Pagination

- 🔼 : Previous Page
- 🔽 : Next Page
