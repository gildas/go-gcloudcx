function uuidv4() {
  return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  )
}

function shortuuid(uuid) {
  // Thanks to: https://github.com/digitalbazaar/base58-universal/blob/master/baseN.js
  const base58 = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
  if (!uuid) uuid = uuidv4()
  const bytes  = new Uint8Array(uuid.toLowerCase().replace(/-/g, '').match(/.{1,2}/g).map(byte => parseInt(byte, 16)))
  const digits = [0]

  let i = 0
  for (i = 0; i < bytes.length; i++) {
      let carry = bytes[i]
      for (let j = 0; j < digits.length; j++) {
          carry += digits[j] << 8
          digits[j] = carry % 58
          carry = (carry / 58) | 0
      }

      while (carry > 0) {
          digits.push(carry % 58)
          carry = (carry / 58) | 0
      }
  }

  let output = ''

  // encoding leading zeros
  for (i = 0; bytes[i] === 0 && i < bytes.length - 1; i++) {
      output += base58[0]
  }

  // encoding the rest
  for (i = digits.length - 1; i >= 0; i--) {
      output += base58[digits[i]]
  }
  return output

  /*
  // When Safari supports BigInt, we can do this:
  let bigint = BigInt('0x' + uuid.toLowerCase().replace(/-/g, ''))

  while (bigint > 0) {
  const mod = Number(bigint % 58n)
  bigint = bigint / 58n
  output.push(base58[mod])
  }
  return output.reverse().join('')
  */
}

const Layout = Vue.component('layout', {
  template: '#layout',
  data() {
    return {
      version: '1.0.0',
    }
  },
})

new Vue({
  el: '#app',
  vuetify: new Vuetify(),
  components: { Layout },
  render(createElement) {
    return createElement(Layout)
  }
})
