# Start Dependencies

```shell
docker-compose up -d
```

# Start Server

```shell
go run main.go
```

# Demo

```shell
curl 'http://localhost:8091/hello?name=world'
# output
{"result":"Hello: world!"}

# Create a cron job
curl -X POST localhost:8091/cronjob \
  -d '{"cronSchedule": "* * * * *"}'

# List all cron jobs
curl 'http://localhost:8091/cronjob?workflowType=CronJobWorkflow'

# Terminate cron job
curl -X DELETE 'localhost:8091/cronjob/my-cron-job'

# Get a cron job
curl -X GET 'localhost:8091/cronjob/my-cron-job'

# Update a cron job schedule
curl -X PUT 'localhost:8091/cronjob' \
  -d '{"cronSchedule":"1 1 * * *", "workflowId":"my-cron-job"}'
```
