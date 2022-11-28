build:
	go build -o dist/gogimage

run:
	dist/gogimage

dev: build run

deploy: build
	scp -pR dist/* root@dpp.st:/var/www/og-img.ld83.org/
	ssh root@dpp.st systemctl restart ogimg