update_unraid:
	@mkdir -p build
	wget https://raw.githubusercontent.com/Squidly271/AppFeed/master/applicationFeed.json -O build/applicationFeed.json

clean:
	rm -rf dist/*

generate:
	mkdir -p dist
	./shctl generate
	./shctl generate -c config.unraid.yml
