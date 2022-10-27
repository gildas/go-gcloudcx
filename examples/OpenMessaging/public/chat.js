const Chat = Vue.component('chat', {
  template: '#chat',
  data() {
    return {
      user:     { userId: 'JohnDoe@KKT', displayName: 'John Doe', pictureUrl: 'https://gravatar.com/avatar/97959eb8244f0cb560e2d30b2075f013?s=400&d=robohash&r=x' },
      chatId:   '',
      iconKey:  0,
      messages: [],
      content:  '',
      socket:   null,
      active:   false,
      dialing:  false,


      fileToSend: [],
      sendingFile: false,
    }
  },
  methods: {
    async startStopChat(event) {
      let self = this
      if (!this.active) {
        console.log("Starting chat", event)
        try {
          let results = await axios.post('/chat', {
            account:    this.account,
            secret:     this.secret,
            webhookUrl: this.webhook,
            userId:     this.user.userId,
          })
          console.log("Received path: ", results.data.path, results)
          this.chatId  = results.data.path.split("/").pop()
          this.dialing = true
          // deepcode ignore MissingClose: close is called in the else branch...
          this.socket  = new WebSocket(`ws://${document.location.host}${results.data.path}`)

          this.socket.onopen = function(event) {
            self.setActive(true)
            self.dialing  = false
            self.messages = []
            self.setfocus('content')
          }
          this.socket.onclose = function(event) {
            console.log('OnClose', event)
            self.setActive(false)
            self.dialing = false
          }
          this.socket.onmessage = this.receiveMessage
        } catch (err) {
          console.error('Failed to connect', err)
        }
      } else {
        console.log("closing chat", event)
        if (this.socket) {
          this.socket.close()
        }
         console.log('the chat is now inactive')
        self.content=""
        self.setActive(false)
      }
    },
    async receiveMessage(event) {
      const lineEvent = JSON.parse(event.data)
      const reqid     = lineEvent.id || uuidv4() // incoming messages contain a UUID not a shortid
      const msgid     = shortuuid(reqid)

      console.log("Received event from LINE Server", event, 'LINE Event: ', lineEvent)
      if (lineEvent.error) {
        console.error('Received Error ', error)
        this.appendTranscript({
          id:        msgid,
          reqid:     reqid,
          from:      'agent',
          type:      'error',
          content:   lineEvent.error,
          timestamp: new Date(),
          tooltip:   event.data,
          //timestamp: event.data.timestamp || new Date(),
        })
        return
      }
      switch (lineEvent.type.toLowerCase()) {
        case 'uploadStatus': {
          const message = this.messages.find(message => message.id === lineEvent.messageId)
          message.status      = 'sending'
          console.log('Sending attachment to LINE Server')
          await this.socket.send(JSON.stringify({
            id:      message.id,
            reqid:   message.reqid,
            userId:  this.user.userId,
            content: lineEvent.content,
          }))
          console.log('Sent to LINE Server')
          message.status      = 'sent'
          message.type        = lineEvent.content.type,
          message.contentType = lineEvent.mimeType
          message.content     = lineEvent.content
          console.log('Updated message', message)
          break
        }
        case "audio": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'audio',
            content:   lineEvent,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        case "file": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'file',
            contentType: lineEvent.contentType, // TODO: Check this!
            content:   lineEvent,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        case "image": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'image',
            content:   lineEvent,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        case "template": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'template',
            content:   lineEvent,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        case "text": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'text',
            content:   lineEvent.text,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        case "video": {
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'video',
            content:   lineEvent,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
          break
        }
        default: {
          console.warn("Not Implemented: ", lineEvent)
          this.appendTranscript({
            id:        msgid,
            reqid:     reqid,
            from:      'agent',
            type:      'text',
            content:   "Not Implemented: \n" + lineEvent.type,
            timestamp: new Date(),
            tooltip:   event.data,
            //timestamp: event.data.timestamp || new Date(),
          })
        }
      }
    },
    async sendMessage() {
      if (!this.content) return
      try {
        const reqid   = uuidv4()
        const shortid = shortuuid(reqid)

        console.log(`Sending message to LINE Server, shortid=${shortid}, reqid=${reqid}`)
        const message = this.appendTranscript({
          id:        shortid,
          reqid:     reqid,
          status:    'sending',
          from:      'guest',
          type:      'text',
          content:   this.content,
          actions:   [],
          timestamp: new Date(),
        })
        await this.socket.send(JSON.stringify({
          id:     shortid,
          reqid:  reqid,
          userId: this.user.userId,
          text:   this.content,
        }))
        message.status = 'sent'

        this.content = ''
        this.setfocus('content')
      } catch (err) {
        console.error('Error while sending message', err)
      }
    },
    async sendLocation() {
      try {
        const reqid    = uuidv4()
        const shortid  = shortuuid(reqid)
        const location = await getLocation()

        console.log(`Sending location to LINE Server, shortid=${shortid}, reqid=${reqid}, location=`, location)
        const message = this.appendTranscript({
          id:          shortid,
          reqid:       reqid,
          status:      'sending',
          from:        'guest',
          type:        'location',
          content:     location,
          actions:     [],
          timestamp:   new Date(),
        })
        await this.socket.send(JSON.stringify({
          id:     shortid,
          reqid:  reqid,
          userId: this.user.userId,
          content: {
            type:        'location',
            id:          shortid,
            title:       'Mock Title',
            address:     'Mock address',
            longitude:   location.coords.longitude,
            latitude:    location.coords.latitude,
          },
        }))
        message.status = 'sent'

        this.content = ''
        this.setfocus('content')
      } catch (err) {
        console.error('Error while sending message', err)
      }
    },
    async sendVideoPlayed({ trackingId }) {
      if (!trackingId) {
        console.log("Cannot send videoPlayed event as trackingId is empty")
        return
      }
      try {
        const reqid   = uuidv4()
        const shortid = shortuuid(reqid)

        console.log(`Sending videoPlayed ${trackingId} to LINE Server, shortid=${shortid}, reqid=${reqid}`)
        await this.socket.send(JSON.stringify({
          id:     shortid,
          reqid:  reqid,
          userId: this.user.userId,
          trackingId,
        }))
      } catch (err) {
        console.error('Error while sending message', err)
      }
    },
    icon(status) {
      switch (status) {
        case 'sent':      return { id: 'mdi-check',     color: 'gray-lighten' }
        case 'displayed': return { id: 'mdi-check-all', color: 'primary' }
        default:          return {}
      }
    },
    setfocus(id) {
      const element = this.$el.querySelector(`#${id}`)

      if (element) {
        this.$nextTick(() => {
          console.log('setting focus to ', element)
          element.focus()
          console.log('done.')
        })
      } else {
        console.error(`Document does not contain any id "${id}"`)
      }
    },
    setActive(state) {
      this.active = state
      this.$emit('chat-active', this.active)
      console.log(`the chat is now ${this.active ? 'active' : 'inactive'}`)
    },
    appendTranscript(message) {
      // TODO: set default stuff, like suggestions: [] and timestamp
      this.messages.push(message)
      this.scrollToEnd()
      return message
    },
    scrollToEnd() {
      /*
      const container = this.$el.querySelector('#chats')
      this.$nextTick(function() {
        container.scrollTop = container.scrollHeight
      })
      */
      this.$nextTick(function() {
        window.scrollTo(0,document.body.scrollHeight)
      })
    },
  },
})
