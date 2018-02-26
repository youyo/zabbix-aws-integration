# zabbix-aws-integration

## Docker image

```
$ docker container run \
	--name zabbix-agent-aws-integration \
	-d \
	-e AWS_ACCESS_KEY_ID=xxxx \
	-e AWS_SECRET_ACCESS_KEY=yyyy \
	-e ZBX_PASSIVESERVERS='0.0.0.0/0' \
	-e ZBX_TIMEOUT=30 \
	-p 10050:10050 \
	youyo/zabbix-aws-integration:latest
$ zabbix_get -s container_host -k zabbix-aws-integration.discovery[ec2,ap-northeast-1,arn:aws:iam::00000000:role/iam_role_name]
```
