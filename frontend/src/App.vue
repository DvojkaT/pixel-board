<template>
  <div class="min-h-full relative min-w-full flex flex-col gap-10 items-center content-center justify-center mt-8">
    <h1 class="font-bold text-2xl">Рисуй</h1>
    <div class="relative">
      <canvas ref="board" class="border border-gray-300"></canvas>
      <canvas ref="overlay" class="absolute inset-0 pointer-events-none"></canvas>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {useCentrifugo} from "@/composables/useCentrifugo.ts";
import {useCanvas} from "@/composables/useCanvas.ts";

const board = ref<HTMLCanvasElement | null>(null);
const overlay = ref<HTMLCanvasElement | null>(null);
const {sub, centrifuge} = useCentrifugo()
useCanvas(board, overlay, sub, centrifuge)
</script>

<style scoped></style>
