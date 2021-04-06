# Image Server

Image Server sets up a simple image server, to let you view folders, images, videos, and compare images among folders.

## How to use

Open your browser and navigate to `http://your_ip:9420` to view your images and videos.

### Browse the current folder

```shell
ImageServer
```

### Browse a specified folder

```shell
ImageServer -root path/to/your/folder
```

### Browse a text file

Here is a `image_list.txt` which contains several image URLs:

```text
http://foo/bar/image1.png
http://foo/bar/image2.jpg
http://foo/bar/image3.gif
```

Then we should run

```shell
ImageServer -root path/to/image_list.txt
```

### Compare folders

The magic parameter is `c`, means __compare__. It follows URL syntax.

Suppose you are now at `dir/a` with the URL of

```text
http://your_ip:9420/dir/a
```

You want to compare this folder with `dir/b`, so you can change the URL to

```text
http://your_ip:9420/dir/a?c=dir/b
```

You then want to compare with `dir_c`, then the url shuld be

```text
http://your_ip:9420/dir/a?c=dir/b&c=dir_c
```
