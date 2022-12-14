<script setup lang="ts">
interface Message {
  username: string
  message: string
}

interface SyncMessage {
  messages: Message[]
}

interface WebsocketMessage<T extends string, D> {
  type: T
  roomId: number
  data: D
}

type WebSocketStatus = 'OPEN' | 'CONNECTING' | 'CLOSED'

const route = useRoute()
const roomId = computed(() => {
  let roomIdStr = ''
  if (Array.isArray(route.params.all))
    roomIdStr = route.params.all[0]

  else
    roomIdStr = route.params.all

  const roomId = Number.parseInt(roomIdStr, 10)
  return isNaN(roomId) ? 0 : roomId
})

const ws = ref<WebSocket | null>()
const messages = ref<Message[]>([])
const wsStatus = ref<WebSocketStatus>('CLOSED')
const chatHistoryEl = ref<HTMLElement>()
const heartbeatInterval = ref()
const syncInterval = ref()

watch(wsStatus, (value) => {
  if (value === 'CLOSED') {
    setTimeout(() => {
      connect()
    }, 10000)
  }
})

function onOpen() {
  wsStatus.value = 'OPEN'
  ws.value?.send(JSON.stringify({
    type: 'sync',
  }))
}

function onMessage(event: MessageEvent) {
  wsStatus.value = 'OPEN'

  try {
    const data = JSON.parse(event.data)
    if (data === '❤️')
      return
    if (data && data.type && data.type === 'message') {
      const messageData = data as WebsocketMessage<'message', Message>
      console.log(roomId.value, messageData.data)
      messages.value.push(messageData.data)
      chatHistoryEl.value?.scrollTo(0, chatHistoryEl.value.scrollHeight)
    }
    if (data && data.type && data.type === 'sync') {
      const messageData = data as WebsocketMessage<'sync', SyncMessage>
      console.log(roomId.value, messageData.data)
      messages.value = messageData.data.messages
      chatHistoryEl.value?.scrollTo(0, chatHistoryEl.value.scrollHeight)
    }
  }
  catch (err) {
    console.error(err)
  }
}

function onClosed() {
  wsStatus.value = 'CLOSED'
  ws.value = undefined
}

function connect() {
  wsStatus.value = 'CONNECTING'
  try {
    ws.value = new WebSocket(`ws://localhost:8123/ws/v1/chat/${roomId.value}`)
  }
  catch (err) {
    console.error(err)
    connect()
  }
}

onMounted(() => {
  connect()
  if (ws.value) {
    ws.value.onopen = onOpen
    ws.value.onmessage = onMessage
    ws.value.onclose = onClosed
    ws.value.onerror = () => {
      connect()
    }

    // syncInterval.value = setInterval(() => {
    //   console.log('syncing')
    //   ws.value.send(JSON.stringify({
    //     type: 'sync',
    //   }))
    // }, 1000)

    heartbeatInterval.value = setInterval(() => {
      if (ws.value)
        ws.value.send('❤️')
    }, 2000)
  }
})

onUnmounted(() => {
  ws.value?.close()
  wsStatus.value = 'CLOSED'

  if (heartbeatInterval.value)
    clearInterval(heartbeatInterval.value)

  if (syncInterval.value)
    clearInterval(syncInterval.value)
})
</script>

<template>
  <div>
    <div class="my-4 rounded-xl bg-blue-200 p-4 text-gray-700 flex">
      <div class="flex-1">
        <span class="font-normal">Chat Room: </span>
        <span class="font-black">{{ roomId }}</span>
      </div>
    </div>
    <div class="p-2 rounded-xl bg-gray-200 w-full">
      <div v-if="!messages || messages.length === 0" class="min-h-[200px] flex items-center">
        <div class="text-center text-gray-400 m-auto">
          No messages yet
        </div>
      </div>
      <div v-else ref="chatHistoryEl" class="">
        <div v-for="(m, index) in messages" :key="index">
          <div class="rounded-xl p-4 bg-white m-4 w-fit shadow-lg">
            {{ m.username }}: {{ m.message }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
