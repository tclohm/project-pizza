# ============ #
# Production
# ============ #

production_host_ip = "54.176.102.78"

# production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh ubuntu@ec2-54-176-102-78.us-west-1.compute.amazonaws.com

## production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api ubuntu@ec2-54-176-102-78.us-west-1.compute.amazonaws.com:~
	rsync -rP --delete ./migrations ubuntu@ec2-54-176-102-78.us-west-1.compute.amazonaws.com:~
	ssh -t ubuntu@ec2-54-176-102-78.us-west-1.compute.amazonaws.com 'migrate -path ~/migrations -database $$PIZZA_DB_DSN up'