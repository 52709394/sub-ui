#!/usr/bin/env bash

pid=$(ps -ef | grep sub-ui | grep -v 'grep' | awk '{print$2}')
if ! [ -z "$pid" ]; then kill $pid; fi
pid=$(ps -ef | grep sub-ui | grep -v 'grep' | awk '{print$2}')
if ! [ -z "$pid" ]; then kill $pid; fi
