<template>
  <div class="min-h-full min-w-full flex flex-col gap-10 items-center content-center justify-center mt-8">
    <h1 class="font-bold text-2xl">Рисуй</h1>
    <canvas ref="board" class="border border-gray-300"></canvas>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {useCentrifugo} from "@/composables/useCentrifugo.ts";
import {useCanvas} from "@/composables/useCanvas.ts";

const board = ref<HTMLCanvasElement | null>(null);
const {sub, centrifuge} = useCentrifugo()
useCanvas(board, sub, centrifuge)
</script>

<style scoped></style>
