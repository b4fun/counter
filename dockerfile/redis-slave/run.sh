#!/bin/sh

redis-server --slaveof redis-service 6379
