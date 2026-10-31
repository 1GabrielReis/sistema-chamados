CREATE TABLE chamados (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    solicitante_id BIGINT NOT NULL,
    responsavel_id BIGINT NOT NULL,
    titulo VARCHAR(150) NOT NULL,
    descricao TEXT NOT NULL,
    prioridade VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'aberto',
    data_abertura TIMESTAMPTZ NOT NULL DEFAULT NOW(),


    CONSTRAINT chamados_solicitante_fk
        FOREIGN KEY (solicitante_id) REFERENCES pessoas(id) ON DELETE RESTRICT,
    CONSTRAINT chamados_responsavel_fk
        FOREIGN KEY (responsavel_id) REFERENCES pessoas(id) ON DELETE RESTRICT,

    CONSTRAINT chamados_titulo_nao_vazio
        CHECK(BTRIM(titulo)<>''),
    CONSTRAINT chamados_descricao_nao_vazia
        CHECK(BTRIM(descricao)<>''),
    CONSTRAINT chamados_prioridade_valida
        CHECK(prioridade  IN ('baixa','media','alta') ),
    CONSTRAINT chamados_status_valido
        CHECK(status  IN ('aberto','em_andamento','resolvido','fechado') )

);

CREATE INDEX chamados_responsavel_status_idx
    ON chamados (responsavel_id, status);

CREATE INDEX chamados_solicitante_idx
    ON chamados (solicitante_id);