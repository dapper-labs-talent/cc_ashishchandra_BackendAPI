#!/bin/sh
set -e

export POSTGRES_USER=dapperlabs
export POSTGRES_PASSWORD=dapperlabs123

# Check if private/public keypair exists
if [ ! -s keys/app-ecdsa ];
then
  echo "`date` Create keypair for the server to sign and verify JWT tokens"
  ssh-keygen -f keys/app-ecdsa -t ecdsa -b 521 -m pem -q -N ""
  RETCODE=$?
  if [ $RETCODE -ne 0 ];
  then
    echo "`date` Could not create the private/public keypair to sign and verify JWT tokens. Please contact the developer. Exiting"
    exit 1
  else
    echo "`date` Successfully created the private/public keypair to sign and verify JWT tokens"
  fi
  ssh-keygen -f keys/app-ecdsa -m pem -e > keys/app-ecdsa.pub
  RETCODE=$?
  if [ $RETCODE -ne 0 ];
  then
    echo "`date` Could not update the public key to PEM format. Please contact the developer. Exiting"
    exit 1
  fi
fi

# Now build docker images
echo "`date` Creating docker images as needed"
docker-compose build
RETCODE=$?

if [ $RETCODE -ne 0 ];
then
  echo "`date` Could not build application image. Please contact the developer. Exiting"
  exit 1
fi

echo "`date` Starting postgres"
docker-compose up -d postgres
RETCODE=$?

if [ $RETCODE -ne 0 ];
then
  echo "`date` Could not start postgres. Exiting"
  exit 1
fi

echo "`date` Waiting for postgres to become ready"
until PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -U $POSTGRES_USER -c '\q'; do
  >&2 echo "`date` Postgres is unavailable - sleeping"
  sleep 1
done

echo "`date` Postgres is ready. Starting application container"
docker-compose up -d server
RETCODE=$?

if [ $RETCODE -ne 0 ];
then
  echo "`date` Could not start server. Please contact the developer. Exiting"
  exit 1
fi
echo
echo
echo "`date` Successfully started all containers. Application is ready on port 3000"
echo
echo
exit 0
