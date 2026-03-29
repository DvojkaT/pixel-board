import {type Centrifuge, Subscription} from "centrifuge";
import {onMounted, ref, type Ref} from "vue";

export function useCanvas(canvasRef: Ref<HTMLCanvasElement | null>, overlayRef: Ref<HTMLCanvasElement | null>, sub: Subscription, centrifuge: Centrifuge) {
  const username = ref<string>("")
  const color = ref<string>("")
  const cursors = new Map<string, {x: number, y: number, name: string, color: string}>()
  const myId = ref<string>("")

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

    const dpr = window.devicePixelRatio || 1

    canvas.width = 1000 * dpr
    canvas.height = 1000 * dpr
    canvas.style.width = '1000px'
    canvas.style.height = '1000px'
    ctx.scale(dpr, dpr)

    canvasOverlay.width = 1000 * dpr
    canvasOverlay.height = 1000 * dpr
    canvasOverlay.style.width = '1000px'
    canvasOverlay.style.height = '1000px'
    overlayCtx.scale(dpr, dpr)

    drawGrid(overlayCtx)

    canvas.addEventListener('mousedown', (e) => {
      const rect = canvas.getBoundingClientRect()
      const x = Math.floor((e.clientX - rect.left) / 10)
      const y = Math.floor((e.clientY - rect.top) / 10)
      sub.publish({ type: 'pixel_paint', x, y, color:color.value })
    })

    let prevX = -1
    let prevY = -1
    let lastSent = 0

    canvas.addEventListener('mousemove', (e) => {
      const now = Date.now()
      if (now - lastSent > 30) {
        const rect = canvas.getBoundingClientRect()
        const mouseX = (e.clientX - rect.left) / 1000
        const mouseY = (e.clientY - rect.top) / 1000
        sub.publish({ type: 'cursor_move', x: mouseX, y: mouseY })

        lastSent = now
      }


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
      if (msg.type === 'cursor_move') {
        if (msg.user_id == myId.value) return
        cursors.set(msg.user_id, { x: msg.x, y: msg.y, name: msg.name, color: msg.color })

        overlayCtx.clearRect(0, 0, 1000, 1000)
        drawGrid(overlayCtx)

        cursors.forEach((cursor) => {

          const cx = cursor.x * 1000
          const cy = cursor.y * 1000

          overlayCtx.fillStyle = cursor.color
          overlayCtx.fillRect(cx - 4, cy - 4, 8, 8)

          overlayCtx.fillStyle = 'black'
          overlayCtx.font = '12px sans-serif'
          overlayCtx.fillText(cursor.name, cx + 6, cy)
        })
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

        username.value = msg.name
        color.value = msg.color
        myId.value = msg.user_id
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

  return {username, color}
}
