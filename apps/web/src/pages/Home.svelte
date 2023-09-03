<script type="ts">
  import { navigate } from "svelte-routing";

  function navigateToTOSpace() {
    navigate("/to-space");
  }

  function navigateToNewGame() {
    navigate("/new-game");
  }
</script>

<div class="container">
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
    <button on:click={navigateToNewGame} class="btn newgame">Novo Jogo</button>
  </div>

  <button on:click={navigateToTOSpace} class="btn to-sess">Espaço T.O</button>
</div>

<style>
  .btn {
    border: none;
    border-radius: 10px;
    color: var(--background);
    font-size: 30px;
    font-weight: 700;
    box-shadow: 0 0 12px 0 rgba(0, 0, 0, 0.24);
    transition: 0.5s;
  }

  .btn:hover {
    box-shadow: 0 0 12px 0 rgba(0, 0, 0, 0.52);
    cursor: pointer;
  }

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

  .container {
    max-width: 800px;
    padding: 20px 20px;
    height: 100%;
    min-height: calc(100vh - 40px);
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
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
