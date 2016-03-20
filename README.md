# docular

Docular is a personal document service which provides a local HTTP server to serve
documents from a given folder (including sub folders).

## Supported file formats

Docular supports the following file formats: HTML, MarkDown, MAFF (Web pages
archived into one single file).

The MarkDown support is from [russross/blackfriday](https://github.com/russross/blackfriday)
with a flavor very similar to github's MarkDown syntax.

## Build docular

Docular is programmed with Go language. At present, docular builds and works fine
on Linux. Other platforms might also work but I didn't have a chance (or desire)
to try it.

If you do not have Go installed, do it now. And don't forget to setup your GOPATH:

```sh
$ export GOPATH = ~/go
```

I use Makefile to build docular, which changes the GOPATH when build.

Next, clone docular source code to a folder out of your GOPATH (for example,
~/tmp):

```sh
$ cd ~/tmp
$ git clone https://github.com/linuxerwang/docular
$ cd docular
```

Now, let's go get the dependent libraries:

```sh
$ make goget
```

If everything goes smooth, you are good to build doculear:

```sh
$ make docular
```

You should be able to find a compiled binary: bin/docular.

If you have my [debmake](https://github.com/linuxerwang/debmaker) installed, you
can easily generate a deb package and then install it:

```sh
$ make deb
$ sudo dpkg -i docular_0.1.0_amd64.deb
```

You can also download the deb package in binaries folder.

## Start docular

When docular is running, it needs certain files which are normally put into folder
/usr/share/docular. If they are not there, copy them there.

Great, now that you have docular installed, let's start it to serve your documents
(suppose your documents are put into ~/mydocs):

```sh
$ docular -doc-dir ~/mydocs
Serving files from /home/zhwang/clients/docular
Server running at http://localhost:3455. CTRL+C to shutdown
```

Now open the URL http://localhost:3455 in your browser. Voila! Everything works!
