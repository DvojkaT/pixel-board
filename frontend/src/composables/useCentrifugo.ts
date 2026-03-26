import {Centrifuge} from "centrifuge";
import {onUnmounted} from "vue";

export  function useCentrifugo() {
  const centrifuge = new Centrifuge("ws://localhost:8080/connection/websocket")
  centrifuge.connect()
  const sub = centrifuge.newSubscription("canvas:main")

  onUnmounted(() => {
    sub.unsubscribe()
    centrifuge.disconnect()
  })

  return {sub, centrifuge}
}
