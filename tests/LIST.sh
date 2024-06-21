#!/bin/sh

bobatool stores list-stores \
	--page_size 3 \
	--insecure \
	--address localhost:8080 \
	--json


bobatool stores list-stores \
	--page_token Mw \
	--page_size 3 \
	--insecure \
	--address localhost:8080 \
	--json
