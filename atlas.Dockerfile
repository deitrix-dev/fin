FROM arigaio/atlas

# Copy the docker binary from the dind image
COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/docker

# Copy the database schema from the fin image
COPY --from=docker.io/deitrix/fin:dev /opt/fin/store/mysql/schema.sql /data/schema.sql