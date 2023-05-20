#!/bin/bash

# Zip current project content
zip -r tk.zip . -x .git\*

# SCP the zipped content to remote server
scp tk.zip mindset_game_tiktok@34.141.121.122:~ 

# with ssh key
# scp tk.zip -i ~/.ssh/gcp/gcp mindset_game_tiktok@34.141.121.122:~ 

# SSH into remote server and execute the following commands
ssh mindset_game_tiktok@34.141.121.122 <<EOF
    cd tk/
    # Stop running docker compose
    docker compose down
    # Remove old project content
    rm -rf ./*
    # Unzip the content of the zipped file into the tk directory
    unzip -o ../tk.zip -d .
    # Build docker compose from the newly deployed content
    docker compose build
    # Start running the newly built docker-compose
    docker compose up -d
EOF
