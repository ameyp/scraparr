#!/bin/bash

CONFIG_FILE=/mounted/config.yml docker run -p 8080:8080 -e CONFIG_FILE -v `pwd`:/mounted ameypar/scraparr
