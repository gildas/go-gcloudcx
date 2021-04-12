const ChatMessage = Vue.component('chat-message', {
  template: '#chat-message',
  props: ['message'],
  data() {
    return {
      calendar:   false,
      dialing:    false,
      replied:    false,
      startTime:   0,
      elapsedTime: 0,
    }
  },
  mounted() {
    console.log('Chat Message', this.message)
    if (this.message.type === 'video') {
      let video = document.getElementById(this.id)

      video.addEventListener('play',    this.onVideoStarted)
      video.addEventListener('playing', this.onVideoResumed)
      video.addEventListener('pause',   this.onVideoPaused)
      video.addEventListener('ended',   this.onVideoEnded)
    }
  },
  computed: {
      classes() {
        return this.message.from === 'guest'
          ? 'd-flex                  speech-bubble speech-bubble-guest pa-2'
          : 'd-flex flex-row-reverse speech-bubble speech-bubble-agent pa-2'
      },
      hasActions() {
        return Array.isArray(message.actions) && message.actions > 0
      },
      id() {
        return `${this.message.type}-${this.message.id}`
      },
      actionId() {
        return 'chat-message-actions-' + this.message.id
      }
  },
  methods: {
    mapUrl(action) {
      return `https://maps.google.com/maps/place/${action.latitude},${action.longitude}`
    },
    onActionClicked(action) {
      console.log('Action event: ', action)
      let displayText = false
      let text = action.displayText || action.text || action.label
      let data = action.data

      this.replied = true
      switch (action.type) {
        case 'message':
          displayText = true
          data = action.text
          break
        case 'URI':
          text = null
          data = null
          break
        case 'location':
          this.$emit('location')
          return
        case 'camera':
          data = action.label
          break
        case 'datetime':
          this.calendar = true
          break
      }
      this.$emit('action', { text, displayText, data })
    },
    onVideoStarted(event) {
      console.log(`Video ${this.id} started`, event)
      this.startTime = new Date().getTime() / 1000
    },
    onVideoResumed(event) {
      console.log(`Video ${this.id} resumed`, event)
      this.startTime = new Date().getTime() / 1000
    },
    onVideoPaused(event) {
      console.log(`Video ${this.id} paused`, event)
      if (this.startTime > 0) {
        let elapsed = new Date().getTime() / 1000 - this.startTime
        this.startTime    = -1
        this.elapsedTime += elapsed
      }
    },
    onVideoEnded(event) {
      console.log(`Video ${this.id} ended`, event)
      let video = event.target
      if (this.startTime > 0) {
        let elapsed = new Date().getTime() / 1000 - this.startTime
        this.startTime    = -1
        this.elapsedTime += elapsed
      }
      console.log(`Elapsed time: ${this.elapsedTime}, duration: ${event.target.duration}`)
      if (this.elapsedTime >= event.target.duration) {
        console.log('Video has been played entirely')
        this.$emit('video-played', { trackingId: this.message.content.trackingId })
      }
    },
  },
})