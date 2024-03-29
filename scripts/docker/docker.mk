include build/docker/.env
export

IMAGES := `docker images --filter "dangling=true" -q --no-trunc`

docker-clean:
	docker rmi ${IMAGES} 2> /dev/null ||:

docker-build:
	docker build -f build/docker/Dockerfile \
      --build-arg VERSION="${VERSION}" \
      -t outdead/rcon .

	docker rmi ${IMAGES} 2> /dev/null ||:

# make docker-run e=pz4 command=players
docker-run:
	docker run -it --rm \
      -v $(CURDIR)/rcon-local.yaml:/rcon.yaml \
      outdead/rcon ./rcon -c rcon.yaml -e $(e) $(command)
