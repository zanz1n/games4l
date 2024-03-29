<script lang="ts">
  import { navigate } from "svelte-routing";
  import { SessionManager, type Session } from "../lib/Session";
  import sharedStyles from "../shared.module.css";
  import { displayFriendlyErr } from "../lib/Repository";
  import { getQuestion, questionsLength, type Question } from "../lib/Question";
  import GameImage from "../components/game/GameImage.svelte";
  import GameSound from "../components/game/GameSound.svelte";
  import { sleep } from "../lib/utils";

  const successVideoEle = document.getElementById(
    "hidden-success"
  )! as HTMLVideoElement;
  const failVideoEle = document.getElementById(
    "hidden-fail"
  ) as HTMLVideoElement;

  interface Info {
    session: Session;
    question: Question;
    questionSpl: string[];
    bigText: boolean;
  }

  const sessionManager = SessionManager.getInstance();

  let disabled = [] as number[];

  async function displaySucess() {
    successVideoEle.classList.remove("hidden-hidden");
    successVideoEle.play();
    await sleep(6000);
    successVideoEle.classList.add("hidden-hidden");
  }

  async function displayFail() {
    failVideoEle.classList.remove("hidden-hidden");
    failVideoEle.play();
    await sleep(4000);
    failVideoEle.classList.add("hidden-hidden");
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
      navigate("/");
      return null;
    }

    const question = getQuestion(session.val.val.currentQuestion);

    if (!question) {
      navigate("/done");
      return null;
    }

    let bigText = false;

    const questionSpl = question.question.split("\n");

    if (question.type == "2Alt") {
      bigText = true;
    } else if (
      question.style != "image" &&
      question.style != "video" &&
      questionSpl.length < 4
    ) {
      bigText = true;
    }

    return {
      session: session.val.val,
      question,
      bigText,
      questionSpl,
    };
  }

  let promise = fetchInfo();

  async function handleSkip() {
    const { session, question } = (await promise) as Info;

    const result = await sessionManager.update(session.name, {
      currentQuestion: session.currentQuestion + 1,
      skiped: session.skiped + 1,
    });

    if (result.err) {
      return alert(displayFriendlyErr(result.val));
    }

    disabled = [];

    promise = fetchInfo();
  }

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
            on:click={function () {
              navigate("/");
            }}
            class={`${sharedStyles.btn} btn-back`}>Voltar</button
          >

          <h1>
            Pergunta {session.session.currentQuestion + 1}/{questionsLength()}
          </h1>
        </div>

        <button
          on:click={function () {
            handleSkip();
          }}
          class={`${sharedStyles.btn} btn-skip`}>Pular</button
        >
      </div>

      <div class="question">
        {#if !!session.question.file}
          {#if session.question.style == "image"}
            <GameImage
              questionType={session.question.type}
              assetName={session.question.file}
            />
          {:else if session.question.style == "audio"}
            <GameSound assetName={session.question.file} />
          {/if}
        {/if}
        <div class="prompt">
          {#each session.questionSpl as line}
            <p class={session.bigText ? "big-p" : "small-p"}>
              {line}
            </p>
          {/each}
        </div>
      </div>

      <div class="answers">
        {#each session.question.answers as answer, i}
          <button
            on:click={function () {
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
    gap: 16px;
  }

  .question .prompt {
    text-align: center;
  }

  .question .prompt .big-p {
    font-size: 28px;
  }

  .question .prompt .small-p {
    font-size: 21px;
  }

  @media screen and (min-height: 880px) {
    .question .prompt .big-p {
      font-size: 30px;
    }

    .question .prompt .small-p {
      font-size: 30px;
    }
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
