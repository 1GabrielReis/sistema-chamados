import { useEffect, useMemo, useState } from "react";
import {
  AlertTriangle,
  ArrowUpRight,
  Check,
  ChevronRight,
  CircleCheckBig,
  CircleDot,
  Clipboard,
  Clock3,
  LayoutDashboard,
  LifeBuoy,
  LoaderCircle,
  Menu,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Sparkles,
  TicketCheck,
  UserPlus,
  UserRoundSearch,
  Users,
  X,
} from "lucide-react";
import { api } from "./api.js";

const STATUS = {
  aberto: "Aberto",
  em_andamento: "Em andamento",
  resolvido: "Resolvido",
  fechado: "Fechado",
};

const PRIORIDADE = {
  baixa: "Baixa",
  media: "Média",
  alta: "Alta",
};

function formatarData(data) {
  if (!data) return "—";
  return new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(data));
}

function App() {
  const [secao, setSecao] = useState("chamados");
  const [menuAberto, setMenuAberto] = useState(false);
  const [chamados, setChamados] = useState([]);
  const [responsaveis, setResponsaveis] = useState([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState("");
  const [modal, setModal] = useState(null);
  const [toast, setToast] = useState(null);

  async function carregarDados(silencioso = false) {
    if (!silencioso) setCarregando(true);
    setErro("");
    try {
      const [listaChamados, listaResponsaveis] = await Promise.all([
        api.listarChamados(),
        api.listarResponsaveis(),
      ]);
      setChamados(listaChamados ?? []);
      setResponsaveis(listaResponsaveis ?? []);
    } catch (error) {
      setErro(error.message || "Não foi possível carregar os dados.");
    } finally {
      setCarregando(false);
    }
  }

  useEffect(() => {
    carregarDados();
  }, []);

  useEffect(() => {
    if (!toast) return undefined;
    const timeout = window.setTimeout(() => setToast(null), 3600);
    return () => window.clearTimeout(timeout);
  }, [toast]);

  function navegar(destino) {
    setSecao(destino);
    setMenuAberto(false);
  }

  function notificar(mensagem, tipo = "sucesso") {
    setToast({ mensagem, tipo });
  }

  async function aoSalvarChamado(mensagem) {
    setModal(null);
    notificar(mensagem);
    await carregarDados(true);
  }

  return (
    <div className="app-shell">
      <Sidebar
        secao={secao}
        menuAberto={menuAberto}
        aoFechar={() => setMenuAberto(false)}
        aoNavegar={navegar}
      />

      <main className="main-content">
        <Topbar
          secao={secao}
          aoAbrirMenu={() => setMenuAberto(true)}
          aoAtualizar={() => carregarDados()}
          carregando={carregando}
        />

        {erro && (
          <div className="connection-alert" role="alert">
            <AlertTriangle size={20} />
            <div>
              <strong>Não foi possível conversar com o servidor.</strong>
              <span>{erro}</span>
            </div>
            <button type="button" className="btn btn--ghost" onClick={() => carregarDados()}>
              Tentar novamente
            </button>
          </div>
        )}

        {secao === "chamados" ? (
          <ChamadosPage
            chamados={chamados}
            carregando={carregando}
            aoCriar={() => setModal({ modo: "criar" })}
            aoVisualizar={(chamado) => setModal({ modo: "visualizar", chamado })}
            aoEditar={(chamado) => setModal({ modo: "editar", chamado })}
          />
        ) : (
          <PessoasPage
            responsaveis={responsaveis}
            aoAtualizar={() => carregarDados(true)}
            aoNotificar={notificar}
          />
        )}
      </main>

      {modal && (
        <ChamadoModal
          key={`${modal.modo}-${modal.chamado?.id ?? "novo"}`}
          modo={modal.modo}
          chamado={modal.chamado}
          responsaveis={responsaveis}
          aoFechar={() => setModal(null)}
          aoEditar={() => setModal({ modo: "editar", chamado: modal.chamado })}
          aoSalvar={aoSalvarChamado}
        />
      )}

      {toast && (
        <div className={`toast toast--${toast.tipo}`} role="status">
          {toast.tipo === "sucesso" ? <Check size={18} /> : <AlertTriangle size={18} />}
          <span>{toast.mensagem}</span>
          <button type="button" aria-label="Fechar aviso" onClick={() => setToast(null)}>
            <X size={16} />
          </button>
        </div>
      )}
    </div>
  );
}

function Sidebar({ secao, menuAberto, aoFechar, aoNavegar }) {
  return (
    <>
      <button
        type="button"
        className={`sidebar-backdrop ${menuAberto ? "is-visible" : ""}`}
        aria-label="Fechar menu"
        onClick={aoFechar}
      />
      <aside className={`sidebar ${menuAberto ? "is-open" : ""}`}>
        <div className="brand">
          <div className="brand__mark">
            <LifeBuoy size={24} strokeWidth={2.2} />
          </div>
          <div>
            <strong>Chamados</strong>
            <span>Central de atendimento</span>
          </div>
        </div>

        <nav className="sidebar__nav" aria-label="Navegação principal">
          <span className="sidebar__label">Workspace</span>
          <button
            type="button"
            className={secao === "chamados" ? "is-active" : ""}
            onClick={() => aoNavegar("chamados")}
          >
            <LayoutDashboard size={19} />
            <span>Visão geral</span>
          </button>
          <button
            type="button"
            className={secao === "pessoas" ? "is-active" : ""}
            onClick={() => aoNavegar("pessoas")}
          >
            <Users size={19} />
            <span>Pessoas</span>
          </button>
        </nav>

        <div className="sidebar__footer">
          <div className="sidebar__status-dot" />
          <div>
            <strong>Ambiente local</strong>
            <span>API em localhost:8080</span>
          </div>
        </div>
      </aside>
    </>
  );
}

function Topbar({ secao, aoAbrirMenu, aoAtualizar, carregando }) {
  return (
    <header className="topbar">
      <div className="topbar__title">
        <button type="button" className="icon-btn mobile-menu" onClick={aoAbrirMenu}>
          <Menu size={21} />
          <span className="sr-only">Abrir menu</span>
        </button>
        <div>
          <span>Central de atendimento</span>
          <strong>{secao === "chamados" ? "Visão geral" : "Pessoas e responsáveis"}</strong>
        </div>
      </div>
      <button type="button" className="btn btn--quiet" onClick={aoAtualizar} disabled={carregando}>
        <RefreshCw size={17} className={carregando ? "is-spinning" : ""} />
        Atualizar
      </button>
    </header>
  );
}

function ChamadosPage({ chamados, carregando, aoCriar, aoVisualizar, aoEditar }) {
  const [busca, setBusca] = useState("");
  const [status, setStatus] = useState("todos");
  const [prioridade, setPrioridade] = useState("todas");

  const filtrados = useMemo(() => {
    const termo = busca.trim().toLowerCase();
    return chamados.filter((chamado) => {
      const correspondeBusca =
        !termo ||
        chamado.titulo.toLowerCase().includes(termo) ||
        chamado.solicitante.nome.toLowerCase().includes(termo) ||
        chamado.responsavel.nome.toLowerCase().includes(termo) ||
        String(chamado.id).includes(termo);
      const correspondeStatus = status === "todos" || chamado.status === status;
      const correspondePrioridade =
        prioridade === "todas" || chamado.prioridade === prioridade;
      return correspondeBusca && correspondeStatus && correspondePrioridade;
    });
  }, [busca, chamados, prioridade, status]);

  const contagem = useMemo(
    () => ({
      total: chamados.length,
      abertos: chamados.filter((item) => item.status === "aberto").length,
      andamento: chamados.filter((item) => item.status === "em_andamento").length,
      concluidos: chamados.filter((item) => ["resolvido", "fechado"].includes(item.status)).length,
    }),
    [chamados],
  );

  return (
    <div className="page-wrap">
      <section className="hero-row">
        <div>
          <span className="eyebrow">Painel operacional</span>
          <h1>Acompanhe cada solicitação</h1>
          <p>Priorize, distribua e resolva chamados em um único lugar.</p>
        </div>
        <button type="button" className="btn btn--primary" onClick={aoCriar}>
          <Plus size={18} />
          Novo chamado
        </button>
      </section>

      <section className="stats-grid" aria-label="Resumo dos chamados">
        <StatCard label="Total de chamados" value={contagem.total} icon={TicketCheck} tone="ink" />
        <StatCard label="Abertos" value={contagem.abertos} icon={CircleDot} tone="blue" />
        <StatCard label="Em andamento" value={contagem.andamento} icon={Clock3} tone="amber" />
        <StatCard label="Concluídos" value={contagem.concluidos} icon={CircleCheckBig} tone="green" />
      </section>

      <section className="panel tickets-panel">
        <div className="panel__header">
          <div>
            <h2>Chamados recentes</h2>
            <p>{filtrados.length} resultado(s) exibido(s)</p>
          </div>
          <div className="filters">
            <label className="search-box">
              <Search size={17} />
              <span className="sr-only">Buscar chamados</span>
              <input
                value={busca}
                onChange={(event) => setBusca(event.target.value)}
                placeholder="Buscar por título, pessoa ou ID"
              />
            </label>
            <label className="compact-select">
              <span className="sr-only">Filtrar por status</span>
              <select value={status} onChange={(event) => setStatus(event.target.value)}>
                <option value="todos">Todos os status</option>
                {Object.entries(STATUS).map(([valor, label]) => (
                  <option key={valor} value={valor}>{label}</option>
                ))}
              </select>
            </label>
            <label className="compact-select">
              <span className="sr-only">Filtrar por prioridade</span>
              <select value={prioridade} onChange={(event) => setPrioridade(event.target.value)}>
                <option value="todas">Todas as prioridades</option>
                {Object.entries(PRIORIDADE).map(([valor, label]) => (
                  <option key={valor} value={valor}>{label}</option>
                ))}
              </select>
            </label>
          </div>
        </div>

        {carregando ? (
          <LoadingState label="Carregando chamados" />
        ) : filtrados.length === 0 ? (
          <EmptyState
            title="Nenhum chamado encontrado"
            description="Ajuste os filtros ou abra uma nova solicitação."
            action="Abrir chamado"
            onAction={aoCriar}
          />
        ) : (
          <>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Chamado</th>
                    <th>Solicitante</th>
                    <th>Responsável</th>
                    <th>Prioridade</th>
                    <th>Status</th>
                    <th>Aberto em</th>
                    <th><span className="sr-only">Ações</span></th>
                  </tr>
                </thead>
                <tbody>
                  {filtrados.map((chamado) => (
                    <tr key={chamado.id} onClick={() => aoVisualizar(chamado)}>
                      <td>
                        <div className="ticket-title">
                          <span>#{String(chamado.id).padStart(4, "0")}</span>
                          <strong>{chamado.titulo}</strong>
                        </div>
                      </td>
                      <td><PersonCell pessoa={chamado.solicitante} /></td>
                      <td><PersonCell pessoa={chamado.responsavel} /></td>
                      <td><PriorityBadge prioridade={chamado.prioridade} /></td>
                      <td><StatusBadge status={chamado.status} /></td>
                      <td className="muted-cell">{formatarData(chamado.data_abertura)}</td>
                      <td>
                        <button
                          type="button"
                          className="icon-btn"
                          aria-label={`Editar chamado ${chamado.id}`}
                          onClick={(event) => {
                            event.stopPropagation();
                            aoEditar(chamado);
                          }}
                        >
                          <Pencil size={16} />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="ticket-cards">
              {filtrados.map((chamado) => (
                <button
                  type="button"
                  className="ticket-card"
                  key={chamado.id}
                  onClick={() => aoVisualizar(chamado)}
                >
                  <div className="ticket-card__top">
                    <span>#{String(chamado.id).padStart(4, "0")}</span>
                    <StatusBadge status={chamado.status} />
                  </div>
                  <strong>{chamado.titulo}</strong>
                  <p>{chamado.descricao}</p>
                  <div className="ticket-card__meta">
                    <PriorityBadge prioridade={chamado.prioridade} />
                    <span>{formatarData(chamado.data_abertura)}</span>
                  </div>
                </button>
              ))}
            </div>
          </>
        )}
      </section>
    </div>
  );
}

function PessoasPage({ responsaveis, aoAtualizar, aoNotificar }) {
  const [nome, setNome] = useState("");
  const [codigoBusca, setCodigoBusca] = useState("");
  const [pessoa, setPessoa] = useState(null);
  const [nomeEdicao, setNomeEdicao] = useState("");
  const [codigoAutorizador, setCodigoAutorizador] = useState("");
  const [salvando, setSalvando] = useState(false);
  const [erroLocal, setErroLocal] = useState("");

  async function cadastrar(event) {
    event.preventDefault();
    setSalvando(true);
    setErroLocal("");
    try {
      const criada = await api.cadastrarPessoa(nome);
      setNome("");
      setPessoa(criada);
      setNomeEdicao(criada.nome);
      setCodigoBusca(criada.registro_publico);
      aoNotificar("Pessoa cadastrada. Guarde o código público.");
    } catch (error) {
      setErroLocal(error.message);
    } finally {
      setSalvando(false);
    }
  }

  async function buscar(event) {
    event.preventDefault();
    setSalvando(true);
    setErroLocal("");
    try {
      const encontrada = await api.buscarPessoa(codigoBusca.trim());
      setPessoa(encontrada);
      setNomeEdicao(encontrada.nome);
    } catch (error) {
      setPessoa(null);
      setErroLocal(error.message);
    } finally {
      setSalvando(false);
    }
  }

  async function atualizar(event) {
    event.preventDefault();
    setSalvando(true);
    setErroLocal("");
    try {
      const atualizada = await api.atualizarPessoa(pessoa.registro_publico, nomeEdicao);
      setPessoa(atualizada);
      aoNotificar("Nome atualizado com sucesso.");
      await aoAtualizar();
    } catch (error) {
      setErroLocal(error.message);
    } finally {
      setSalvando(false);
    }
  }

  async function habilitar(event) {
    event.preventDefault();
    setSalvando(true);
    setErroLocal("");
    try {
      const atualizada = await api.habilitarResponsavel(
        pessoa.registro_publico,
        codigoAutorizador.trim(),
      );
      setPessoa(atualizada);
      setCodigoAutorizador("");
      aoNotificar("Pessoa habilitada como responsável.");
      await aoAtualizar();
    } catch (error) {
      setErroLocal(error.message);
    } finally {
      setSalvando(false);
    }
  }

  async function copiarCodigo(codigo) {
    try {
      await navigator.clipboard.writeText(codigo);
      aoNotificar("Código copiado para a área de transferência.");
    } catch {
      aoNotificar("Não foi possível copiar o código.", "erro");
    }
  }

  return (
    <div className="page-wrap">
      <section className="hero-row">
        <div>
          <span className="eyebrow">Cadastros e permissões</span>
          <h1>Pessoas e responsáveis</h1>
          <p>Cadastre solicitantes, consulte códigos e gerencie responsáveis.</p>
        </div>
      </section>

      {erroLocal && (
        <div className="inline-error" role="alert">
          <AlertTriangle size={18} />
          {erroLocal}
        </div>
      )}

      <section className="people-grid">
        <article className="panel form-card">
          <div className="card-heading">
            <span className="card-heading__icon"><UserPlus size={20} /></span>
            <div>
              <h2>Cadastrar pessoa</h2>
              <p>O código público será criado automaticamente.</p>
            </div>
          </div>
          <form onSubmit={cadastrar}>
            <label className="field">
              <span>Nome completo</span>
              <input
                required
                value={nome}
                onChange={(event) => setNome(event.target.value)}
                placeholder="Ex.: Ana Souza"
              />
            </label>
            <button type="submit" className="btn btn--primary btn--full" disabled={salvando}>
              {salvando ? <LoaderCircle className="is-spinning" size={18} /> : <UserPlus size={18} />}
              Cadastrar pessoa
            </button>
          </form>
        </article>

        <article className="panel form-card">
          <div className="card-heading">
            <span className="card-heading__icon card-heading__icon--blue">
              <UserRoundSearch size={20} />
            </span>
            <div>
              <h2>Localizar pessoa</h2>
              <p>Use o código público recebido no cadastro.</p>
            </div>
          </div>
          <form onSubmit={buscar}>
            <label className="field">
              <span>Código público</span>
              <div className="input-with-action">
                <input
                  required
                  value={codigoBusca}
                  onChange={(event) => setCodigoBusca(event.target.value.toUpperCase())}
                  placeholder="Cole o código aqui"
                />
                <button type="submit" className="icon-btn" aria-label="Buscar pessoa">
                  <Search size={18} />
                </button>
              </div>
            </label>
          </form>
        </article>
      </section>

      {pessoa && (
        <section className="panel person-profile">
          <div className="person-profile__summary">
            <Avatar nome={pessoa.nome} large />
            <div>
              <div className="person-profile__name-line">
                <h2>{pessoa.nome}</h2>
                {pessoa.eh_responsavel && (
                  <span className="responsible-chip"><ShieldCheck size={14} /> Responsável</span>
                )}
              </div>
              <button
                type="button"
                className="public-code"
                onClick={() => copiarCodigo(pessoa.registro_publico)}
              >
                {pessoa.registro_publico}
                <Clipboard size={14} />
              </button>
            </div>
          </div>

          <div className="person-profile__actions">
            <form onSubmit={atualizar} className="inline-form">
              <label className="field">
                <span>Atualizar nome</span>
                <input
                  required
                  value={nomeEdicao}
                  onChange={(event) => setNomeEdicao(event.target.value)}
                />
              </label>
              <button type="submit" className="btn btn--secondary" disabled={salvando}>
                <Pencil size={16} /> Salvar nome
              </button>
            </form>

            {!pessoa.eh_responsavel && (
              <form onSubmit={habilitar} className="inline-form inline-form--accent">
                <label className="field">
                  <span>Código do responsável autorizador</span>
                  <input
                    required
                    value={codigoAutorizador}
                    onChange={(event) => setCodigoAutorizador(event.target.value.toUpperCase())}
                    placeholder="Ex.: RESP-001"
                  />
                </label>
                <button type="submit" className="btn btn--primary" disabled={salvando}>
                  <ShieldCheck size={16} /> Habilitar
                </button>
              </form>
            )}
          </div>
        </section>
      )}

      <section className="panel responsibles-panel">
        <div className="panel__header">
          <div>
            <h2>Responsáveis ativos</h2>
            <p>{responsaveis.length} pessoa(s) disponível(is) para distribuição</p>
          </div>
          <span className="responsible-count"><ShieldCheck size={16} /> {responsaveis.length}</span>
        </div>
        <div className="responsibles-grid">
          {responsaveis.map((responsavel) => (
            <article className="responsible-card" key={responsavel.registro_publico}>
              <Avatar nome={responsavel.nome} />
              <div>
                <strong>{responsavel.nome}</strong>
                <button type="button" onClick={() => copiarCodigo(responsavel.registro_publico)}>
                  {responsavel.registro_publico} <Clipboard size={13} />
                </button>
              </div>
              <span className="availability-dot" title="Disponível" />
            </article>
          ))}
        </div>
      </section>
    </div>
  );
}

function ChamadoModal({ modo, chamado, responsaveis, aoFechar, aoEditar, aoSalvar }) {
  const visualizando = modo === "visualizar";
  const editando = modo === "editar";
  const [salvando, setSalvando] = useState(false);
  const [erro, setErro] = useState("");
  const [form, setForm] = useState({
    codigo_solicitante: chamado?.solicitante.registro_publico ?? "",
    codigo_responsavel: chamado?.responsavel.registro_publico ?? "",
    distribuicao_automatica: false,
    titulo: chamado?.titulo ?? "",
    descricao: chamado?.descricao ?? "",
    prioridade: chamado?.prioridade ?? "media",
    status: chamado?.status ?? "aberto",
  });

  function alterar(campo, valor) {
    setForm((atual) => ({ ...atual, [campo]: valor }));
  }

  async function enviar(event) {
    event.preventDefault();
    setSalvando(true);
    setErro("");
    try {
      if (editando) {
        await api.atualizarChamado(chamado.id, {
          codigo_responsavel: form.distribuicao_automatica ? "" : form.codigo_responsavel,
          distribuicao_automatica: form.distribuicao_automatica,
          titulo: form.titulo,
          descricao: form.descricao,
          prioridade: form.prioridade,
          status: form.status,
        });
        await aoSalvar("Chamado atualizado com sucesso.");
      } else {
        await api.cadastrarChamado({
          codigo_solicitante: form.codigo_solicitante,
          codigo_responsavel: form.distribuicao_automatica ? "" : form.codigo_responsavel,
          distribuicao_automatica: form.distribuicao_automatica,
          titulo: form.titulo,
          descricao: form.descricao,
          prioridade: form.prioridade,
        });
        await aoSalvar("Chamado aberto com sucesso.");
      }
    } catch (error) {
      setErro(error.message);
    } finally {
      setSalvando(false);
    }
  }

  return (
    <div className="modal-layer" role="presentation" onMouseDown={aoFechar}>
      <section
        className={`modal ${visualizando ? "modal--details" : ""}`}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <header className="modal__header">
          <div>
            <span className="eyebrow">
              {visualizando ? `Chamado #${String(chamado.id).padStart(4, "0")}` : editando ? "Editar chamado" : "Nova solicitação"}
            </span>
            <h2 id="modal-title">
              {visualizando ? chamado.titulo : editando ? "Atualize os dados" : "Abrir novo chamado"}
            </h2>
          </div>
          <button type="button" className="icon-btn" onClick={aoFechar} aria-label="Fechar">
            <X size={20} />
          </button>
        </header>

        {visualizando ? (
          <TicketDetails chamado={chamado} aoEditar={aoEditar} />
        ) : (
          <form className="ticket-form" onSubmit={enviar}>
            {erro && <div className="inline-error"><AlertTriangle size={17} /> {erro}</div>}

            {!editando && (
              <label className="field field--full">
                <span>Código do solicitante</span>
                <input
                  required
                  value={form.codigo_solicitante}
                  onChange={(event) => alterar("codigo_solicitante", event.target.value.toUpperCase())}
                  placeholder="Código público recebido no cadastro"
                />
              </label>
            )}

            <label className="field field--full">
              <span>Título do chamado</span>
              <input
                required
                maxLength={150}
                value={form.titulo}
                onChange={(event) => alterar("titulo", event.target.value)}
                placeholder="Resuma o problema em uma frase"
              />
            </label>

            <label className="field field--full">
              <span>Descrição</span>
              <textarea
                required
                rows={5}
                value={form.descricao}
                onChange={(event) => alterar("descricao", event.target.value)}
                placeholder="Descreva o que aconteceu e como podemos ajudar"
              />
            </label>

            <label className="field">
              <span>Prioridade</span>
              <select value={form.prioridade} onChange={(event) => alterar("prioridade", event.target.value)}>
                {Object.entries(PRIORIDADE).map(([valor, label]) => (
                  <option key={valor} value={valor}>{label}</option>
                ))}
              </select>
            </label>

            {editando && (
              <label className="field">
                <span>Status</span>
                <select value={form.status} onChange={(event) => alterar("status", event.target.value)}>
                  {Object.entries(STATUS).map(([valor, label]) => (
                    <option key={valor} value={valor}>{label}</option>
                  ))}
                </select>
              </label>
            )}

            <div className="assignment-card field--full">
              <div className="assignment-card__head">
                <span><Sparkles size={18} /></span>
                <div>
                  <strong>Distribuição automática</strong>
                  <p>Encaminha para quem possui menos chamados abertos.</p>
                </div>
                <label className="switch">
                  <input
                    type="checkbox"
                    checked={form.distribuicao_automatica}
                    onChange={(event) => alterar("distribuicao_automatica", event.target.checked)}
                  />
                  <span />
                </label>
              </div>

              {!form.distribuicao_automatica && (
                <label className="field">
                  <span>Responsável</span>
                  <select
                    required
                    value={form.codigo_responsavel}
                    onChange={(event) => alterar("codigo_responsavel", event.target.value)}
                  >
                    <option value="">Selecione um responsável</option>
                    {responsaveis.map((item) => (
                      <option key={item.registro_publico} value={item.registro_publico}>
                        {item.nome} · {item.registro_publico}
                      </option>
                    ))}
                  </select>
                </label>
              )}
            </div>

            <footer className="modal__footer field--full">
              <button type="button" className="btn btn--quiet" onClick={aoFechar}>Cancelar</button>
              <button type="submit" className="btn btn--primary" disabled={salvando}>
                {salvando ? <LoaderCircle size={18} className="is-spinning" /> : editando ? <Check size={18} /> : <Plus size={18} />}
                {salvando ? "Salvando..." : editando ? "Salvar alterações" : "Abrir chamado"}
              </button>
            </footer>
          </form>
        )}
      </section>
    </div>
  );
}

function TicketDetails({ chamado, aoEditar }) {
  return (
    <div className="ticket-details">
      <div className="ticket-details__badges">
        <StatusBadge status={chamado.status} />
        <PriorityBadge prioridade={chamado.prioridade} />
      </div>

      <p className="ticket-details__description">{chamado.descricao}</p>

      <div className="people-pair">
        <div>
          <span>Solicitante</span>
          <PersonCell pessoa={chamado.solicitante} showCode />
        </div>
        <ChevronRight size={20} />
        <div>
          <span>Responsável</span>
          <PersonCell pessoa={chamado.responsavel} showCode />
        </div>
      </div>

      <dl className="details-list">
        <div><dt>Identificador</dt><dd>#{String(chamado.id).padStart(4, "0")}</dd></div>
        <div><dt>Aberto em</dt><dd>{formatarData(chamado.data_abertura)}</dd></div>
      </dl>

      <footer className="modal__footer">
        <button type="button" className="btn btn--primary" onClick={aoEditar}>
          <Pencil size={17} /> Editar chamado
        </button>
      </footer>
    </div>
  );
}

function StatCard({ label, value, icon: Icon, tone }) {
  return (
    <article className="stat-card">
      <span className={`stat-card__icon stat-card__icon--${tone}`}><Icon size={20} /></span>
      <div><span>{label}</span><strong>{value}</strong></div>
      <ArrowUpRight size={17} className="stat-card__arrow" />
    </article>
  );
}

function StatusBadge({ status }) {
  return <span className={`badge badge--status badge--${status}`}>{STATUS[status] ?? status}</span>;
}

function PriorityBadge({ prioridade }) {
  return (
    <span className={`badge badge--priority badge--priority-${prioridade}`}>
      <span /> {PRIORIDADE[prioridade] ?? prioridade}
    </span>
  );
}

function PersonCell({ pessoa, showCode = false }) {
  return (
    <div className="person-cell">
      <Avatar nome={pessoa.nome} />
      <div>
        <strong>{pessoa.nome}</strong>
        <span>{showCode ? pessoa.registro_publico : pessoa.eh_responsavel ? "Responsável" : "Solicitante"}</span>
      </div>
    </div>
  );
}

function Avatar({ nome, large = false }) {
  const iniciais = nome
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((parte) => parte[0])
    .join("")
    .toUpperCase();
  return <span className={`avatar ${large ? "avatar--large" : ""}`}>{iniciais || "?"}</span>;
}

function LoadingState({ label }) {
  return (
    <div className="loading-state">
      <LoaderCircle className="is-spinning" size={28} />
      <span>{label}</span>
    </div>
  );
}

function EmptyState({ title, description, action, onAction }) {
  return (
    <div className="empty-state">
      <span><TicketCheck size={28} /></span>
      <h3>{title}</h3>
      <p>{description}</p>
      <button type="button" className="btn btn--secondary" onClick={onAction}>
        <Plus size={17} /> {action}
      </button>
    </div>
  );
}

export default App;
