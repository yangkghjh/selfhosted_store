update_unraid:
	wget https://raw.githubusercontent.com/Squidly271/AppFeed/master/applicationFeed.json -O build/applicationFeed.json

clean:
	rm -rf dist/*