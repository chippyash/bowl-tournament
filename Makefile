VERSION:=`cat VERSION`
NAME=bowl-tournament
DESC="Utility to a bowls tournament graph of games"

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: clean  ## Build the binary and zip into dist/bowl-tournament.tar.gz
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/bowl-tournament bowl-tournament.go
	chmod a+x dist/bowl-tournament
	cp dist/bowl-tournament .
	cp dependencies/plantuml.jar .
	tar -czvf dist/bowl-tournament.tar.gz bowl-tournament README.md plantuml.jar
	rm -f bowl-tournament
	rm -f plantuml.jar

.PHONY: deploy-local
deploy-local: build  ## Deploy binary locally
	sudo mv dist/bowl-tournament /usr/bin/bowl-tournament
	sudo mkdir -p /usr/share/bowl-tournament/
	sudo cp dist/plantuml.jar /usr/share/bowl-tournament/plantuml.jar

.PHONY: clean
clean: ## Clean up the dist directory
	rm -f ./dist/bowl-tournament*

.PHONY: getfpm
getfpm: ## Install FPM and supporting libraries
	sudo yum install ruby-devel gcc rpm-build rubygems -y
	sudo gem install --no-document fpm

.PHONY: package
package: build ## Package executable into rpm and deb packages
	fpm -t rpm -v $(VERSION) --rpm-use-file-permissions --force
	fpm -t deb -v $(VERSION) --rpm-use-file-permissions --force
	mv -f *.rpm ./dist
	mv -f *.deb ./dist
	rm -f dist/bowl-tournament
