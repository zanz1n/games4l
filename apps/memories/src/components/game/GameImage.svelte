<script lang="ts">
  export let assetName: string;
  export let questionType: "2Alt" | "4Alt";

  function resolveAssetAddress(name: string) {
    return "/questions/images/" + name;
  }

  const imageSrc = resolveAssetAddress(assetName) + ".png";
</script>

<div class="image-box">
  <img
    class={questionType == "2Alt" ? "big-img" : "small-img"}
    src={imageSrc}
    alt={assetName}
  />
</div>

<style>
  .image-box .small-img {
    max-height: 240px;
  }

  .image-box .big-img {
    max-height: 320px;
  }

  @media screen and (max-height: 880px) {
    .image-box .small-img {
      max-height: 150px;
    }

    .image-box .big-img {
      max-height: 240px;
    }
  }

  @media screen and (max-height: 660px) {
    .image-box .small-img {
      max-height: 120px;
    }

    .image-box .big-img {
      max-height: 150px;
    }
  }
</style>
