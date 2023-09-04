<script lang="ts">
  import { navigate } from "svelte-routing";
  import { SessionManager, type Session } from "../lib/Session";
  import sharedStyles from "../shared.module.css";
  import { displayFriendlyErr } from "../lib/Repository";
  import { getQuestion, questionsLength, type Question } from "../lib/Question";

  interface Info {
    session: Session;
    question: Question;
  }

  const sessionManager = SessionManager.getInstance();

  let disabled = [] as number[];

  async function displaySucess() {
    confirm("Correto");
  }

  async function displayFail() {
    confirm("Tente novamente");
  }

  async function fetchInfo(): Promise<Info | null> {
    const session = await sessionManager.getCurrent();

    if (session.err) {
      alert(displayFriendlyErr(session.val));
      navigate("/");
      return null;
    }

    if (session.val.none) {
      alert("Selecione ou crie uma conta antes de iniciar o jogo!");
      navigate("/") as never;
      return null;
    }

    const question = getQuestion(session.val.val.currentQuestion);

    if (!question) {
      navigate("/done");
      return null;
    }

    return {
      session: session.val.val,
      question,
    };
  }

  let promise = fetchInfo();

  async function handleSubmission(n: number) {
    const { session, question } = (await promise) as Info;

    if (question.correct_answer == n) {
      const result = await sessionManager.update(session.name, {
        currentQuestion: session.currentQuestion + 1,
        hits: session.hits + 1,
      });

      if (result.err) {
        return alert(displayFriendlyErr(result.val));
      }

      disabled = [];

      await displaySucess();
    } else {
      const result = await sessionManager.update(session.name, {
        errors: session.errors + 1,
      });

      if (result.err) {
        return alert(displayFriendlyErr(result.val));
      }

      disabled = [...disabled, n];

      console.log(disabled);

      await displayFail();
    }

    promise = fetchInfo();
  }
</script>

<div class={sharedStyles.container}>
  {#await promise then session}
    {#if !!session}
      <div class="top">
        <div class="left">
          <button
            on:click={() => {
              navigate("/");
            }}
            class={`${sharedStyles.btn} btn-back`}>Voltar</button
          >

          <h1>
            Pergunta {session.session.currentQuestion + 1}/{questionsLength()}
          </h1>
        </div>

        <button class={`${sharedStyles.btn} btn-skip`}>Pular</button>
      </div>

      <div class="question">
        <p>{session.question.question}</p>
      </div>

      <div class="answers">
        {#each session.question.answers as answer, i}
          <button
            on:click={() => {
              handleSubmission(i);
            }}
            disabled={disabled.includes(i)}
            class={`${sharedStyles.btn} answer-btn answer-btn${i}`}
            >{answer}</button
          >
        {/each}
      </div>
    {:else}
      <p>Algo deu errado!</p>
    {/if}
  {:catch err}
    <p>Algo deu errado: {err}</p>
  {/await}
</div>

<style>
  .top {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }

  .top .left {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 20px;
  }

  .top .left h1 {
    font-size: 22px;
  }

  .btn-back {
    background-color: var(--cyan);
    font-size: 24px;
    width: 100px;
    height: 36px;
  }

  .btn-skip {
    background-color: var(--cyan);
    font-size: 24px;
    width: 100px;
    height: 36px;
  }

  .question {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .question p {
    font-size: 28px;
  }

  .answers {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
  }

  .answer-btn {
    width: 85%;
    height: 64px;
    color: var(--foreground);
    font-size: 26px;
  }

  /* .answer-btn:disabled {
    background-color: var(--foreground);
  } */

  .answer-btn:disabled:hover {
    cursor: not-allowed;
  }

  .answer-btn0 {
    background-color: var(--red);
  }
  .answer-btn0:disabled {
    background-color: var(--red-disabled);
  }
  .answer-btn1 {
    background-color: var(--yellow);
  }
  .answer-btn1:disabled {
    background-color: var(--yellow-disabled);
  }
  .answer-btn2 {
    background-color: var(--purple);
  }
  .answer-btn2:disabled {
    background-color: var(--purple-disabled);
  }
  .answer-btn3 {
    background-color: var(--cyan);
  }
  .answer-btn3:disabled {
    background-color: var(--cyan-disabled);
  }
</style>
