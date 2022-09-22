# ============ #
# Production
# ============ #

production_host_ip = "54.176.102.78"

# production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh $AWS

## production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api $AWS
	rsync -rP --delete ./migrations $AWS
	ssh -t $AWS 'sudo migrate -path ~/migrations -database $$PIZZA_DB_DSN up'