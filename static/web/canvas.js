const centrifuge = new Centrifuge("ws://localhost:8080/connection/websocket")
const sub = centrifuge.newSubscription("canvas:main")

sub.subscribe()
centrifuge.connect()

const canvas = document.getElementById('board')
const ctx = canvas.getContext('2d')
canvas.width = 1000
canvas.height = 1000

function drawPixel(x, y, color) {
    ctx.fillStyle = '#' + color.toString(16).padStart(6, '0')
    ctx.fillRect(x * 10, y * 10, 10, 10)
}

canvas.addEventListener('click', (e) => {
    const rect = canvas.getBoundingClientRect()
    const x = Math.floor((e.clientX - rect.left) / 10)
    const y = Math.floor((e.clientY - rect.top) / 10)
    sub.publish({ type: 'pixel_paint', x, y, color: 0xFF0000 })
})

sub.on('publication', (ctx) => {
    const msg = ctx.data
    if (msg.type === 'pixel_paint') {
        drawPixel(msg.x, msg.y, msg.color)
    }
})