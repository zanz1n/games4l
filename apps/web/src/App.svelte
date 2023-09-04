<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";
  import NewGame from "./pages/NewGame.svelte";
  import Game from "./pages/Game.svelte";

  export let url = "";
</script>

<Router {url}>
  <Route path="/" component={Home} />
  <Route path="/new-game" component={NewGame} />
  <Route path="/game" component={Game} />
</Router>
