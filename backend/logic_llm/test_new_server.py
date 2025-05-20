import pytest
from fastapi.testclient import TestClient
from logic_llm.server import app
import os
from dotenv import load_dotenv,get_key
from datetime import datetime, timedelta, timezone
from jose import jwt 

if not load_dotenv("../.env"):
    print("Could not load .env file")
    exit(1)

client = TestClient(app,"http://127.0.0.1:3000")

def create_jwt_token():
    # Définir les revendications du token
    payload = {
        "sub": "test_user",
        "exp": datetime.now(tz=timezone.utc) + timedelta(hours=1),  # Expiration dans 1 heure
        "iat": datetime.now(tz=timezone.utc),  # Heure d'émission
        "scope": "user"
    }
    # Générer le token JWT
    token = jwt.encode(payload, os.getenv("ACCESS_SECRET_KEY_LLM"), algorithm="HS256")
    return token

def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}

def test_receive_data():
    # Remplacez par un token valide pour vos tests
    # token = "Bearer " + os.getenv("ACCESS_SECRET_KEY_LLM")
    token = create_jwt_token()
    data = {
        "user_id": "test_user",
        "conversation_id": "test_convo",
        "message": "Hello, how are you?"
    }
    headers = {"Authorization": f"Bearer {token}"}
    response = client.post("/llm-protected", json=data, headers=headers)
    assert response.status_code == 200
    # Vous pouvez ajouter plus d'assertions ici pour vérifier le contenu de la réponse

def test_receive_data_empty_message():
    # token = "Bearer " + os.getenv("ACCESS_SECRET_KEY_LLM")
    token = create_jwt_token()
    data = {
        "user_id": "test_user",
        "conversation_id": "test_convo",
        "message": ""
    }
    headers = {"Authorization": f"Bearer {token}"}
    response = client.post("/llm-protected", json=data, headers=headers)
    assert response.status_code == 200
    # Vérifiez que le message d'erreur est retourné
    assert "user message empty" in response.text