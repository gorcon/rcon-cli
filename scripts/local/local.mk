# If the first argument is "docker-deploy"...
ifeq (compile, $(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

compile:
	sh scripts/local/compile.sh $(RUN_ARGS)

run:
	sh scripts/local/run.sh

lint:
	golangci-lint run
