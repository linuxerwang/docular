export GOPATH := ${GOPATH}:$(shell pwd)
export VERSION := "0.1.0"

all:
	@echo "make goget   : Fetch the dependent packages"
	@echo "make docular : build docular"
	@echo "make deb     : build deb package"

goget:
	go get github.com/microcosm-cc/bluemonday
	go get github.com/shurcooL/github_flavored_markdown
	go get github.com/go-martini/martini
	go get github.com/microcosm-cc/bluemonday
	go get github.com/russross/blackfriday
	go get github.com/sergi/go-diff/diffmatchpatch
	go get github.com/shurcooL/github_flavored_markdown
	go get github.com/shurcooL/highlight_diff
	go get github.com/shurcooL/highlight_go
	go get github.com/shurcooL/sanitized_anchor_name
	go get github.com/sourcegraph/annotate
	go get github.com/sourcegraph/syntaxhighlight
	go get golang.org/x/crypto/bcrypt
	go get golang.org/x/net

docular:
	go install -ldflags "-s -w" docular/cmd/docular

deb:
	debmaker -version ${VERSION} -spec-file deb-docular.spec
