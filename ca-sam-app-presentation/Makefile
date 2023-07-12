.PHONY: build

build:
	sam build

local:
	sam local start-api --docker-network go-academy --port 3000

invoke:
	sam local invoke ScheduleFunction