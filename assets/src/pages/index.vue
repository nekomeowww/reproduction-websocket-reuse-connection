<script setup lang="ts">
const roomId = ref('')

const router = useRouter()
const go = () => {
  if (roomId.value && !isNaN(Number.parseInt(roomId.value, 10)))
    router.push(`/room/${encodeURIComponent(roomId.value)}`)
}
</script>

<template>
  <div class="flex flex-col items-center">
    <div text-4xl>
      <div i-carbon-campsite inline-block />
    </div>

    <div py-4 />

    <TheInput
      v-model="roomId"
      placeholder="Enter a room ID"
      autocomplete="false"
      @keydown.enter="go"
    />

    <div>
      <button
        btn m-3 text-sm
        :disabled="!roomId"
        @click="go"
      >
        Enter
      </button>
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: home
</route>
