import {type Centrifuge, Subscription} from "centrifuge";
import {onMounted, type Ref} from "vue";

export function useCanvas(canvasRef: Ref<HTMLCanvasElement | null>, overlayRef: Ref<HTMLCanvasElement | null>, sub: Subscription, centrifuge: Centrifuge) {
  onMounted(() => {
    sub.subscribe()
    const canvas = canvasRef.value
    if (!canvas) {
      return
    }

    const canvasOverlay = overlayRef.value
    if (!canvasOverlay) {
      return
    }

    const ctx = canvas.getContext('2d')!
    if (!ctx) return

    const overlayCtx = canvasOverlay.getContext('2d')!
    if (!overlayCtx) return

    canvas.width = 1000
    canvas.height = 1000

    canvasOverlay.width = 1000
    canvasOverlay.height = 1000

    drawGrid(overlayCtx)

    canvas.addEventListener('mousedown', (e) => {
      const rect = canvas.getBoundingClientRect()
      const x = Math.floor((e.clientX - rect.left) / 10)
      const y = Math.floor((e.clientY - rect.top) / 10)
      sub.publish({ type: 'pixel_paint', x, y, color:'#FF0000' })
    })

    let prevX = -1
    let prevY = -1
    canvas.addEventListener('mousemove', (e) => {
      const rect = canvas.getBoundingClientRect()
      const x = Math.floor((e.clientX - rect.left) / 10)
      const y = Math.floor((e.clientY - rect.top) / 10)

      if (x === prevX && y === prevY) return

      if (prevX >= 0) {
        overlayCtx.clearRect(0, 0, 1000, 1000) //todo Убийца оптимизации. Может сделать 3-ий канвас на highlight?
        drawGrid(overlayCtx)
      }

      overlayCtx.fillStyle = 'rgba(136,120,120,0.4)'
      overlayCtx.fillRect(x * 10, y * 10, 10, 10)

      prevX = x
      prevY = y
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

    centrifuge.on('message', (ctx) => {
      const msg = ctx.data
      if (msg.type === 'canvas_state') {
        for (let x = 0; x < msg.pixels.length; x++) {
          for (let y = 0; y < msg.pixels[x].length; y++) {
            if (msg.pixels[x][y] !== '') {
              drawPixel(x, y, msg.pixels[x][y])
            }
          }
        }
      }
    })
  })

  function drawGrid(overlayCtx: CanvasRenderingContext2D) {
    overlayCtx.strokeStyle = 'rgba(0, 0, 0, 0.15)'
    overlayCtx.lineWidth = 0.5

    overlayCtx.beginPath()

    for (let i = 0; i <= 1000; i += 10) {
      overlayCtx.moveTo(i, 0)
      overlayCtx.lineTo(i, 1000)
      overlayCtx.moveTo(0, i)
      overlayCtx.lineTo(1000, i)
    }

    overlayCtx.stroke()
  }
}
