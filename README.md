

Configure AWS CLI to work with LocalStack:
```shell
~/.aws/config
[profile localstack]
region=us-east-1
output=json
endpoint_url = http://localhost:4566

~/.aws/credentials
[localstack]
aws_access_key_id=test
aws_secret_access_key=test
```

Create new secret in LocalStack:
```shell
aws secretsmanager create-secret \
    --name "postgres-db-secret-1" \
    --secret-string '{"username": "dbuser1", "password": "dbpass1", "host": "localhost", "port": 5432}' --profile localstack

aws secretsmanager list-secrets --profile localstack
```
