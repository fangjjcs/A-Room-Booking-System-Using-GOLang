#!/bin/bash

# -o output to a file called bookings
go build -o bookings cmd/web/*.go && ./bookings