<script lang="ts">
  import { navigate } from "svelte-routing";
  import sharedStyles from "../shared.module.css";

  function navigateToTOSpace() {
    navigate("/to-space");
  }

  function navigateToNewGame() {
    navigate("/new-game");
  }
</script>

<div class={sharedStyles.container}>
  <div class="header">
    <img
      height="50px"
      width="150px"
      src="/hospiral-reger-logo.png"
      alt="Hospital Reger"
    />
    <img
      height="40px"
      width="130px"
      src="/poliedro-logo.png"
      alt="Hospital Reger"
    />
  </div>

  <div class="main">
    <img
      width="280px"
      height="280px"
      src="/memories-logo.png"
      alt="Memories For Life Logo"
    />
    <button on:click={navigateToNewGame} class={`${sharedStyles.btn} newgame`}
      >Novo Jogo</button
    >
  </div>

  <button on:click={navigateToTOSpace} class={`${sharedStyles.btn} to-sess`}
    >Espaço T.O</button
  >
</div>

<style>
  .newgame {
    background-color: var(--cyan);
    width: 256px;
    height: 68px;
    font-size: 38px;
  }

  .to-sess {
    width: 170px;
    height: 48px;
    background-color: var(--purple);
    font-size: 26px;
  }

  .header {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }

  .main {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
  }
</style>
