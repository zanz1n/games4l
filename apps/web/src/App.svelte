<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/home.svelte";
  import NewGame from "./pages/new-game.svelte";
  import Game from "./pages/game.svelte";
  import Done from "./pages/done.svelte";

  export let url = "";
</script>

<Router {url}>
  <Route path="/" component={Home} />
  <Route path="/new-game" component={NewGame} />
  <Route path="/game" component={Game} />
  <Route path="/done" component={Done} />
</Router>
