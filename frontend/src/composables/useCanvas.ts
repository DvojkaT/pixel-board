import {Subscription} from "centrifuge";
import {onMounted, type Ref} from "vue";

export function useCanvas(canvasRef: Ref<HTMLCanvasElement | null>, sub: Subscription) {
  onMounted(() => {
    sub.subscribe()
    const canvas = canvasRef.value
    if (!canvas) {
      return
    }
    const ctx = canvas.getContext('2d')!
    if (!ctx) return

    canvas.width = 1000
    canvas.height = 1000

    canvas.addEventListener('click', (e) => {
      const rect = canvas.getBoundingClientRect()
      const x = Math.floor((e.clientX - rect.left) / 10)
      const y = Math.floor((e.clientY - rect.top) / 10)
      sub.publish({ type: 'pixel_paint', x, y, color:'#FF0000' })
    })

    function drawPixel(x: number, y: number, color: string) {
      ctx.fillStyle = color
      ctx.fillRect(x * 10, y * 10, 10, 10)
    }

    sub.on('publication', (ctx) => {
      const msg = ctx.data
      if (msg.type === 'pixel_paint') {
        drawPixel(msg.x, msg.y, msg.color)
      }
    })
  })
}
