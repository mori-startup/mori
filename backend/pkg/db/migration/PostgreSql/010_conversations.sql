-- +migrate Up
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY, -- Identifiant unique de la conversation
    user_id VARCHAR(100) NOT NULL, -- Identifiant de l'utilisateur, clé étrangère vers la table users
    conversation_id VARCHAR(100) NOT NULL, -- Identifiant unique de la conversation
    convo JSONB[] NOT NULL, -- Liste des conversations
    new_conversation BOOLEAN DEFAULT FALSE NOT NULL, -- Indique si c'est une nouvelle conversation
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS conversations;
