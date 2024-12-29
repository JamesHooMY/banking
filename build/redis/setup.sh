#!/bin/sh

PORT=7000
NODES=6
REPLICAS=${REPLICAS:-1}
HOST_IP=$(hostname -i)
REDIS_PASSWORD=${REDIS_PASSWORD:-redis_password}

echo "Waiting for Redis nodes to start..."
sleep 5

echo "Starting cluster setup..."
echo "HOST_IP is $HOST_IP"

HOSTS=""
ENDPORT=7005
CURRENTPORT=$PORT

while [ $CURRENTPORT -le $ENDPORT ]; do
    HOSTS="$HOSTS $HOST_IP:$CURRENTPORT"
    CURRENTPORT=$((CURRENTPORT+1))
done

echo "Creating cluster with hosts: $HOSTS"
redis-cli -a "$REDIS_PASSWORD" --cluster create $HOSTS --cluster-replicas $REPLICAS

tail -f /dev/null