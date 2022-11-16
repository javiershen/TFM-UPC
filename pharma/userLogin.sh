#!/bin/bash
#
#
# Exit on first error


echo USERNAME_ID:
read username
export USERNAME_ID="$username"
echo USER_PASSWORD:
read userpassword
export USER_PASSWORD="$password"
echo Your user ID is "${username}"