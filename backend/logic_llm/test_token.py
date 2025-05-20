import pytest
from fastapi.testclient import TestClient
from logic_llm.server import app, verify_token
from fastapi import HTTPException, Depends
from fastapi.security import HTTPAuthorizationCredentials


client = TestClient(app)

# Test de la fonction verify_token
def test_verify_token_valid():
    # Créer un jeton JWT valide
    from jose import jwt 
    import os
    from datetime import datetime, timedelta,timezone
    from dotenv import load_dotenv
    
    load_dotenv()

    # Clé secrète utilisée pour signer le jeton
    secret_key = os.getenv("ACCESS_SECRET_KEY_LLM")
    
    # Payload du jeton
    payload = {
        "user_id": "12345",
        "exp": datetime.now(timezone.utc) + timedelta(minutes=30)
    }

    # Générer un jeton
    token = jwt.encode(payload, secret_key, algorithm="HS256")

    credentials = HTTPAuthorizationCredentials(scheme="Bearer", credentials=token)
    
    # Appeler la fonction verify_token
    try:
        result = verify_token(credentials)
        assert result["user_id"] == "12345"
    except HTTPException:
        pytest.fail("verify_token a levé une HTTPException avec un jeton valide.")

def test_verify_token_expired():
    # Créer un jeton JWT expiré
    from jose import jwt 
    import os
    from datetime import datetime, timedelta,timezone

    secret_key = os.getenv("ACCESS_SECRET_KEY_LLM")

    payload = {
        "user_id": "12345",
        "exp": datetime.now(timezone.utc) - timedelta(minutes=7)  # Jeton expiré
    }

    token = jwt.encode(payload, secret_key, algorithm="HS256")
    credentials = HTTPAuthorizationCredentials(scheme="Bearer", credentials=token)

    with pytest.raises(HTTPException) as excinfo:
        verify_token(credentials)

    assert excinfo.value.status_code == 401
    assert "JWT Error: Signature has expired." == str(excinfo.value.detail)

def test_verify_token_invalid():
    # Jeton invalide
    invalid_token = "invalid_token"
    credentials = HTTPAuthorizationCredentials(scheme="Bearer", credentials=invalid_token)

    with pytest.raises(HTTPException) as excinfo:
        verify_token(credentials)

    assert excinfo.value.status_code == 401
    assert "JWT Error: Not enough segments" == str(excinfo.value.detail)
