
set -e

sh /wait-for-db.sh

./migrate up

exec "$@"
