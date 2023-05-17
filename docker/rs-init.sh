#!/bin/bash

mongosh <<EOF
var config = {
    "_id": "dbrs",
    "version": 1,
    "members": [
        {
            "_id": 1,
            "host": "mongo_primary:27017",
            "priority": 3
        },
        {
            "_id": 2,
            "host": "mongo_secondary:27017",
            "priority": 2
        },
        {
            "_id": 3,
            "host": "mongo_tertiary:27017",
            "priority": 1
        }
    ]
};
rs.initiate(config, { force: true });
rs.status();
EOF
