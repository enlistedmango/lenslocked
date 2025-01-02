#!/bin/sh
/app/goose -dir /app/migrations postgres "$DATABASE_URL" up 
