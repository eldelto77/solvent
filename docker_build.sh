#!/bin/bash
mkdir docker_images
docker build --tag="solvent:latest" .
docker save -o ./docker_images/solvent_image.tar solvent
