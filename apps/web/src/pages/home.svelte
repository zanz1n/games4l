<script lang="ts">
  // @ts-ignore
  import HospitalRegerLogo from "../assets/hospiral-reger-logo.png?w=225px&format=avif;webp;png&as=picture";
  // @ts-ignore
  import MemoriesForLifeLogo from "../assets/memories-logo.png?w=420px&format=avif;webp;png&as=picture";
  // @ts-ignore
  import PoliedroLogo from "../assets/poliedro-logo.png?w=195px&format=avif;webp;png&as=picture";

  import { navigate } from "svelte-routing";
  import sharedStyles from "../shared.module.css";
  import Picture from "../components/Picture.svelte";

  function navigateToTOSpace() {
    navigate("/to-space");
  }

  function navigateToNewGame() {
    navigate("/new-game");
  }
</script>

<div class={sharedStyles.container}>
  <div class="header">
    <Picture meta={HospitalRegerLogo} sizes="150px" alt="Hospital Reger" />
    <Picture meta={PoliedroLogo} sizes="130px" alt="Poliedro" />
  </div>

  <div class="main">
    <Picture
      loading="eager"
      meta={MemoriesForLifeLogo}
      sizes="280px"
      alt="Memories for Life"
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
