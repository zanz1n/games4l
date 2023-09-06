<script lang="ts">
  // @ts-ignore
  import ToScreenPic from "../assets/to-screen.png?w=480px&format=avif;webp;png&as=picture";

  import { navigate } from "svelte-routing";
  import Picture from "../components/Picture.svelte";
  import sharedStyles from "../shared.module.css";
  import { SessionManager, type Session } from "../lib/Session";
  import { displayFriendlyErr } from "../lib/Repository";
  import { questionsLength } from "../lib/Question";

  interface Info {
    sessions: Session[];
    totalQuestions: number;
  }

  function formatDate(millis: number): string {
    const date = new Date(millis);

    return (
      `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()} ` +
      `${date.getHours()}:${date.getMinutes()}`
    );
  }

  const sessionManager = SessionManager.getInstance();

  let current: Session | null = null;

  async function handleContinue() {
    if (!current) return;

    const result = await sessionManager.setCurrent(current.name);
    if (result.err) {
      return alert(displayFriendlyErr(result.val));
    }

    if (current.currentQuestion >= questionsLength()) {
      return alert(`${current.name} já completou todas as questões!`);
    }

    navigate("/game");
  }

  async function getData(): Promise<Info | null> {
    const session = await sessionManager.getAll();

    if (session.err) {
      alert(displayFriendlyErr(session.val));
      navigate("/");
      return null;
    }

    if (session.val.length > 0) {
      current = session.val[0];
    }

    return {
      sessions: session.val,
      totalQuestions: questionsLength(),
    };
  }

  let promise = getData();

  async function setCurrent(name: string) {
    const result = await sessionManager.get(name);

    if (result.err) {
      return alert(displayFriendlyErr(result.val));
    }

    if (result.val.none) {
      alert("Algo deu errado! Não foi possível encontrar o usuário!");
      promise = getData();
      return;
    }

    current = result.val.val;
  }

  async function deleteSession(name: string) {
    const cf = confirm("Deseja mesmo excluir " + name + "?");

    if (cf) {
      const result = await sessionManager.delete(name);

      if (result.err) {
        return alert(displayFriendlyErr(result.val));
      }

      promise = getData();
    }
  }
</script>

<div class={`${sharedStyles.container} container`}>
  <div class="top">
    <button
      on:click={function () {
        navigate("/");
      }}
      class={`${sharedStyles.btn} btn-back`}>Voltar</button
    >
    <Picture
      loading="eager"
      meta={ToScreenPic}
      sizes="320px"
      alt="Tela da To"
    />
    <h1>Pacientes</h1>
  </div>

  {#await promise then info}
    {#if !!info}
      <div class={info.sessions.length == 0 ? "" : "mid"}>
        {#if info.sessions.length == 0}
          <p class="big-p">Não há nenhum paciente por enquanto!</p>
        {:else}
          <div class="left">
            {#each info.sessions as session}
              <div class="session">
                <h2>{session.name}</h2>
                <div class="btn-row">
                  <button
                    on:click={function () {
                      setCurrent(session.name);
                    }}
                    class={`${sharedStyles.btn} btn-info`}>Ver Info</button
                  >

                  <button
                    on:click={function () {
                      deleteSession(session.name);
                    }}
                    class={`${sharedStyles.btn} btn-del`}>Excluir</button
                  >
                </div>
              </div>
            {/each}
          </div>

          <div class="right">
            {#if !!current}
              <div class="info">
                <div class="field">
                  <h2>Nome:</h2>
                  <p>{current.name}</p>
                </div>
                <hr />
                <div class="field">
                  <h2>Idade:</h2>
                  <p>{current.age}</p>
                </div>
                <hr />
                <div class="field small">
                  <h2>Começou:</h2>
                  <p>{formatDate(current.createdAt)}</p>
                </div>
                <hr />
                <div class="field">
                  <h2>Questão:</h2>
                  <p>{current.currentQuestion}/{info.totalQuestions}</p>
                </div>
                <hr />
                <div class="field">
                  <h2>Acertos:</h2>
                  <p>{current.hits}/{info.totalQuestions}</p>
                </div>
                <hr />
                <div class="field">
                  <h2>Erros:</h2>
                  <p>{current.errors}</p>
                </div>
                <hr />
                <div class="field">
                  <h2>Puladas:</h2>
                  <p>{current.skiped}/{info.totalQuestions}</p>
                </div>
              </div>

              <div class="bottom">
                <button
                  on:click={function () {
                    handleContinue();
                  }}
                  disabled={current.currentQuestion == info.totalQuestions}
                  class={`${sharedStyles.btn} btn-continue`}>Continuar</button
                >
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}
  {/await}

  <div />
</div>

<style>
  .container {
    gap: 32px;
  }
  .btn-back {
    background-color: var(--cyan);
    font-size: 24px;
    width: 100px;
    height: 36px;
  }

  .top {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
  }

  .top button {
    align-self: flex-start;
  }

  .big-p {
    text-align: center;
    font-size: 28px;
  }

  .mid {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-self: center;
    align-items: center;
    width: 100%;
  }

  .mid .left {
    display: flex;
    flex-direction: column;
    gap: 20px;
    width: 60%;
    margin: 10px;
  }

  .left .session {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    gap: 32px;
    padding: 12px;
    background-color: #e7e7e7;
    border-radius: 12px;
    box-shadow: 0 0 12px 0 rgba(0, 0, 0, 0.24);
    transition: 0.5s;
  }

  .left .session:hover {
    box-shadow: 0 0 12px 0 rgba(0, 0, 0, 0.42);
  }

  .session .btn-row {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .btn-row button {
    font-size: 20px;
    padding: 5px 8px;
  }

  .btn-row .btn-del {
    background-color: var(--red);
  }

  .btn-row .btn-info {
    background-color: var(--purple);
  }

  .mid .right {
    margin: 10px;
    width: 30%;
  }

  .right .info {
    display: flex;
    flex-direction: column;
  }

  .info .field {
    display: flex;
    flex-direction: row;
    align-items: flex-end;
    justify-content: space-between;
    gap: 16px;
  }

  .field p {
    font-size: 18px;
    margin-bottom: 2px;
  }

  .small p {
    font-size: 15px;
  }

  .field h2 {
    font-size: 22px;
  }

  .right .bottom {
    margin-top: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .right .bottom .btn-continue {
    padding: 6px;
    background-color: var(--green);
    align-self: center;
  }

  .right .bottom .btn-continue:disabled {
    background-color: var(--green-disabled);
  }

  .right .bottom .btn-continue:disabled:hover {
    cursor: not-allowed;
  }

  @media screen and (max-width: 660px) {
    .mid {
      flex-direction: column-reverse;
      gap: 32px;
    }

    .mid .left {
      width: 90%;
    }

    .mid .right {
      width: 90%;
    }
  }
</style>
