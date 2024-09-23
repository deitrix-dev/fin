FROM arigaio/atlas

# Copy the dind binary from the dind image
COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/docker