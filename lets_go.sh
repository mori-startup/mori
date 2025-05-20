#!/bin/bash

# Démarrer le serveur Go
echo "Démarrage du serveur Go..."
cd backend
go run server.go &
BACKEND_GO_PID=$!

# Démarrer Uvicorn
echo "Démarrage de Uvicorn..."
uvicorn server:app --host 127.0.0.1 --port 3000 &
UVICORN_PID=$!

# Démarrer le frontend
echo "Démarrage du frontend..."
cd ../frontend
npm run serve &
FRONTEND_PID=$!

echo "Tous les serveurs sont démarrés."

# Optionnel : Attendre que tous les processus se terminent
wait