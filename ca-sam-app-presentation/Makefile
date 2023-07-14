.PHONY: build

build:
	sam build

local:
	sam local start-api --docker-network go-academy-presentation --port 3000 --env-vars env.json

invoke:
	sam local invoke ScheduleFunction