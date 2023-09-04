<script lang="ts">
  // @ts-ignore
  import NewGameImage from "../assets/new-game.png?w=350px&format=avif;webp;png&as=picture";

  import { navigate } from "svelte-routing";
  import sharedStyles from "../shared.module.css";
  import Picture from "../components/Picture.svelte";
  import { SessionManager } from "../lib/Session";
  import { displayFriendlyErr } from "../lib/Repository";

  const sessionManager = SessionManager.getInstance();

  function onSubmit() {
    const name = (document.getElementById("name") as HTMLInputElement).value;
    if (1 > name.length) {
      return alert("Insira um nome válido!");
    }

    const age = (document.getElementById("age") as HTMLInputElement)
      .valueAsNumber;
    if (isNaN(age) || 1 > age) {
      return alert("Insira uma idade válida!");
    }

    sessionManager.setCurrent(name).then(async (result) => {
      if (result.err) {
        return alert(displayFriendlyErr(result.val));
      }

      result = await sessionManager.create(name, age);

      if (result.err) {
        return alert(displayFriendlyErr(result.val));
      }

      navigate("/game");
    });
  }
</script>

<div class={sharedStyles.container}>
  <div class="center top">
    <div class="header">
      <button
        on:click={() => {
          navigate("/");
        }}
        class={`${sharedStyles.btn} btn-back`}>Voltar</button
      >
    </div>
    <Picture meta={NewGameImage} sizes="350px" alt="New Game" />
  </div>

  <div class="main center">
    <div class="title">
      <h1>Identificação do paciente</h1>
    </div>

    <input name="name" id="name" type="text" placeholder="Nome" />
    <input name="age" id="age" type="number" placeholder="Idade" />
  </div>

  <div class="center">
    <button on:click={onSubmit} class={`${sharedStyles.btn} btn-start`}
      >Prosseguir</button
    >
  </div>
</div>

<style>
  .btn-back {
    background-color: var(--cyan);
    font-size: 24px;
    width: 100px;
    height: 36px;
  }

  .btn-start {
    background-color: var(--cyan);
    width: 70%;
    min-width: 280px;
    height: 68px;
    font-size: 38px;
    margin-top: 24px;
  }

  .center {
    display: flex;
    align-items: center;
    flex-direction: column;
  }

  .top {
    gap: 32px;
  }

  .top .header {
    width: 100%;
  }

  .main {
    gap: 16px;
  }

  .main .title {
    text-align: center;
  }

  .main input {
    height: 52px;
    width: 60%;
    min-width: 280px;
    font-size: 24px;
    padding: 0 8px;
  }

  .main .title {
    margin-bottom: 24px;
  }
</style>
