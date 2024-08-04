FROM loadimpact/k6

COPY ./scripts /scripts
COPY ./config.json /config.json

# No need to create /results directory in Dockerfile
# ENTRYPOINT ["k6", "run", "--config", "/config.json", "--out", "json=/results/output.json", "/scripts/get-users.js"]
ENTRYPOINT ["k6", "run", "--config", "/config.json", "--summary-export", "/results/output.json", "/scripts/get-users.js"]
