all: build

BUILDARG=-ldflags " -s -X main.buildtime=`date '+%Y-%m-%d_%H:%M:%S'` -X main.githash=`git rev-parse HEAD`"

build:
	go build ${BUILDARG}

clean:
	rm -f mcrontab
