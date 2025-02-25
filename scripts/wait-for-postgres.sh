set -e

until PGPASSWORD=$DB_PASSWORD psql -h "db" -U "postgres" -d "swift_codes" -c '\q'; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - executing command"
exec "$@"