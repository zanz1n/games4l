<script lang="ts">
  export let assetName: string;

  function resolveAssetAddress(name: string) {
    return "/questions/audio/" + name;
  }

  const audioSrc = resolveAssetAddress(assetName) + ".";
</script>

<audio controls>
  <source type="audio/webm" src={audioSrc + "webm"} />
  <source type="audio/mpeg" src={audioSrc + "mp3"} />
  <p class="err-msg">Seu navegador não suporta áudio!</p>
</audio>

<style>
  .err-msg {
    color: var(--red);
    font-size: 18px;
    font-weight: 700;
  }
</style>
