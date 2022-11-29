build:
	go build -o dist/gogimage

run:
	dist/gogimage

dev: build run

deploy: build
	rsync -arvz dist/* root@dpp.st:/var/www/og-img.ld83.com/
	ssh root@dpp.st systemctl restart ogimg