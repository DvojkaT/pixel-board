<template>
  <div
    class="min-h-full relative min-w-full flex flex-col gap-10 items-center content-center justify-center mt-8">
    <div class="grid grid-cols-2">
      <h1 class="font-bold text-2xl">Рисуй</h1>
      <div class="flex flex-col">
        <div>
          <span>Ваше имя: {{ username }}</span>
        </div>
        <div class="flex items-center gap-2">
          <span>Ваш цвет:</span>
          <div :style="{ backgroundColor: color }" class="w-5 h-5 rounded"></div>
        </div>
      </div>
    </div>
    <div class="relative">
      <canvas ref="board" class="border border-gray-300"></canvas>
      <canvas ref="overlay" class="border border-gray-300 absolute inset-0 pointer-events-none"></canvas>
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
const {username, color} = useCanvas(board, overlay, sub, centrifuge)
</script>

<style scoped></style>
