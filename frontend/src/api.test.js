import assert from "node:assert/strict";
import { afterEach, test } from "node:test";
import { api } from "./api.js";

const fetchOriginal = globalThis.fetch;

afterEach(() => {
  globalThis.fetch = fetchOriginal;
});

test("listarChamados consulta o endpoint esperado", async () => {
  globalThis.fetch = async (url, options) => {
    assert.equal(url, "/api/chamados");
    assert.equal(options.headers.Accept, "application/json");

    return {
      ok: true,
      status: 200,
      json: async () => [{ id: 1, titulo: "Impressora" }],
    };
  };

  const chamados = await api.listarChamados();

  assert.equal(chamados.length, 1);
  assert.equal(chamados[0].titulo, "Impressora");
});

test("listarSolicitantes consulta o novo endpoint", async () => {
  globalThis.fetch = async (url, options) => {
    assert.equal(url, "/api/pessoas/solicitantes");
    assert.equal(options.headers.Accept, "application/json");

    return {
      ok: true,
      status: 200,
      json: async () => [
        {
          registro_publico: "PESSOA-001",
          nome: "Pessoa solicitante",
          eh_responsavel: false,
        },
      ],
    };
  };

  const solicitantes = await api.listarSolicitantes();

  assert.equal(solicitantes.length, 1);
  assert.equal(solicitantes[0].registro_publico, "PESSOA-001");
  assert.equal(solicitantes[0].eh_responsavel, false);
});

test("cadastrarChamado envia JSON pelo método POST", async () => {
  const entrada = {
    codigo_solicitante: "PESSOA-001",
    codigo_responsavel: "RESP-001",
    distribuicao_automatica: false,
    titulo: "Computador travando",
    descricao: "O computador trava ao iniciar.",
    prioridade: "alta",
  };

  globalThis.fetch = async (url, options) => {
    assert.equal(url, "/api/chamados");
    assert.equal(options.method, "POST");
    assert.deepEqual(JSON.parse(options.body), entrada);

    return {
      ok: true,
      status: 201,
      json: async () => ({ id: 10, ...entrada }),
    };
  };

  const chamado = await api.cadastrarChamado(entrada);

  assert.equal(chamado.id, 10);
});

test("a API transforma resposta de erro em exceção legível", async () => {
  globalThis.fetch = async () => ({
    ok: false,
    status: 400,
    text: async () => "prioridade inválida\n",
  });

  await assert.rejects(
    () => api.cadastrarChamado({}),
    /prioridade inválida/,
  );
});
