<script lang="ts">
  import FireworkGif from "../assets/firework-animation1.gif";

  import { navigate } from "svelte-routing";
  import { questionsLength } from "../lib/Question";
  import { displayFriendlyErr } from "../lib/Repository";
  import { SessionManager } from "../lib/Session";
  import sharedStyles from "../shared.module.css";

  const questionLen = questionsLength();
  const sessionManager = SessionManager.getInstance();

  const promise = (async () => {
    const session = await sessionManager.getCurrent();
    if (session.err) {
      return alert(displayFriendlyErr(session.val));
    }
    if (session.val.none) {
      return navigate("/");
    }

    const { val } = session.val;

    if (questionLen == 0 || val.currentQuestion < questionLen) {
      return navigate("/game");
    }

    const result = await sessionManager.setCurrent(null!);
    if (result.err) {
      return alert(displayFriendlyErr(result.val));
    }

    setTimeout(() => {
      navigate("/");
    }, 7000);
  })();
</script>

<div
  class={sharedStyles.container}
  style={`background-image: url(${FireworkGif});`}
>
  {#await promise then}
    <button
      on:click={() => {
        navigate("/");
      }}
      class={`${sharedStyles.btn} btn-back`}>Voltar</button
    >

    <div class="mid">
      <h1>Parabéns!</h1>
    </div>

    <button
      on:click={() => {
        navigate("/to-space");
      }}
      class={`${sharedStyles.btn} to-sess`}>Espaço T.O</button
    >
  {/await}
</div>

<style>
  .btn-back {
    background-color: var(--cyan);
    font-size: 24px;
    width: 100px;
    height: 36px;
  }

  .mid h1 {
    font-size: 72px;
  }

  .mid {
    align-self: center;
  }

  .to-sess {
    width: 170px;
    height: 48px;
    background-color: var(--purple);
    font-size: 26px;
  }
</style>
