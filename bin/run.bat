COPY .env.example .env
docker-compose --file docker-compose.yaml --env-file .env up --build