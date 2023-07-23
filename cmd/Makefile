.PHONY: build

build:
	sam build --no-cached

local:
	sam local start-api --warm-containers eager --docker-network go-academy-presentation --port 3000 --env-vars env.json

invoke:
	sam local invoke ScheduleFunction --docker-network go-academy-presentation --env-vars env.json