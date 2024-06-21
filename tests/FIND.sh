#!/bin/sh

bobatool stores find-stores \
	--bounds.min.latitude 28.08 \
	--bounds.max.latitude 48.09 \
	--bounds.min.longitude -96.9 \
	--bounds.max.longitude -86.8 \
	--insecure \
	--address localhost:8080 \
	--json


