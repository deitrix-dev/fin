#!/usr/bin/env bash

script_path=$(realpath $0)
script_dir=$(dirname $script_path)

atlas schema inspect --env prod --format '{{ sql . }}' > "$script_dir/../store/mysql/schema.sql"