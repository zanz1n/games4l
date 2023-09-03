<script lang="ts">
  let counter = 1;

  function setCounter(n: number) {
    counter = n;
  }
</script>

<main>
  <h1>Hello World</h1>
  <p>Counter: {counter}</p>
  <button
    on:click={() => {
      setCounter(counter - 1);
    }}>Decrement</button
  >
  <button
    on:click={() => {
      setCounter(counter + 1);
    }}>Increment</button
  >
</main>

<style>
</style>
