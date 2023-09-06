<script lang="ts">
  import { resolveAssetAddress } from "../../lib/Question";

  export let assetName: string;

  const audioSrc = resolveAssetAddress(assetName) + ".webm";
</script>

<audio controls >
  <source src={audioSrc} />
  <p class="err-msg">Seu navegador não suporta áudio!</p>
</audio>

<style>
  .err-msg {
    color: var(--red);
    font-size: 18px;
    font-weight: 700;
  }
</style>
