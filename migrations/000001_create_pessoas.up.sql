CREATE TABLE pessoas (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    registro_publico TEXT NOT NULL UNIQUE,
    nome VARCHAR(100) NOT NULL,
    eh_responsavel BOOLEAN NOT NULL DEFAULT FALSE,
    data_cadastro TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    

    -- Restrições 
    CONSTRAINT pessoas_registro_publico_nao_vazio
        CHECK (BTRIM(registro_publico) <> ''),

    CONSTRAINT pessoas_nome_nao_vazio
        CHECK(BTRIM(nome) <> '')
    );