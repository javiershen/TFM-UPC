#!/bin/bash
#
#
# Exit on first error


echo ENTITY_ID pharmacy2 hospital1:
read entity
export ENTITY_ID="$entity"
echo USERNAME sanitaryUser adminLab pharmacyUser pharmacyAdmin:
read username
export USERNAME="$username"
echo USER_PASSWORD adminpw:
read userpassword
export USER_PASSWORD="$password"
echo You are logged with user "${USERNAME}"
