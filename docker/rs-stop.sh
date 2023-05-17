#!/bin/bash

mongosh <<EOF
db.adminCommand({ shutdown: 1, force: true });
EOF
