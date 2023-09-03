<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";

  export let url = "";
</script>

<Router {url}>
  <Route path="/" component={Home} />
</Router>
