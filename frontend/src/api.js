const API_BASE = import.meta.env?.VITE_API_URL || "/api";

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      Accept: "application/json",
      ...(options.body ? { "Content-Type": "application/json" } : {}),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const message = (await response.text()).trim();
    throw new Error(message || `Falha na requisição (${response.status})`);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

export const api = {
  listarChamados: () => request("/chamados"),
  buscarChamado: (id) => request(`/chamados/${id}`),
  cadastrarChamado: (dados) =>
    request("/chamados", {
      method: "POST",
      body: JSON.stringify(dados),
    }),
  atualizarChamado: (id, dados) =>
    request(`/chamados/${id}`, {
      method: "PUT",
      body: JSON.stringify(dados),
    }),
  cadastrarPessoa: (nome) =>
    request("/pessoas", {
      method: "POST",
      body: JSON.stringify({ nome }),
    }),
  buscarPessoa: (codigo) => request(`/pessoas/${encodeURIComponent(codigo)}`),
  atualizarPessoa: (codigo, nome) =>
    request(`/pessoas/${encodeURIComponent(codigo)}`, {
      method: "PATCH",
      body: JSON.stringify({ nome }),
    }),
  listarResponsaveis: () => request("/pessoas/responsaveis"),
  habilitarResponsavel: (codigoPessoa, codigoAutorizador) =>
    request(`/pessoas/${encodeURIComponent(codigoPessoa)}/responsavel`, {
      method: "PATCH",
      body: JSON.stringify({ codigo_autorizador: codigoAutorizador }),
    }),
};
