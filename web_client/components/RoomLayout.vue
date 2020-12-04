<template>
  <div id="room-layout" />
</template>

<script>
import common from '@/plugins/common'
export default {
  name: 'RoomLayout',
  props: [],
  data: () => ({
    layoutFile: '@/assets/english-room.html',
    emptySeatColor: '#ddcec3',
    maxDisplayNameLength: 6,
  }),
  async created() {
    const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_layout')
    url.search = new URLSearchParams({
      room_id: 'english-room',
    }).toString()
    const response = await fetch(url.toString(), { method: 'GET' })
    const res_html = await response.text()
    const dom_parser = new DOMParser()
    try {
      const layout_doc = dom_parser.parseFromString(this.escape_html(res_html), 'text/html')
      document.getElementById('room-layout').appendChild(layout_doc.body)
      console.log('どう？')
      console.log(layout_doc)
    } catch (e) {
      console.error(e)
    }
  },
  methods: {
    escape_html(string) {
      if (typeof string !== 'string') {
        return string
      }
      return string.replace(/[&'`"<>]/g, function (match) {
        return {
          '&': '&amp;',
          "'": '&#x27;',
          '`': '&#x60;',
          '"': '&quot;',
          '<': '&lt;',
          '>': '&gt;',
        }[match]
      })
    },
    window: (onload = function () {
      const numSeats = 23

      for (let n = 0; n < numSeats; n++) {
        const seatNum = n + 1
        const seatElement = document.getElementById('seat-' + seatNum.toString())
        if (seatElement) {
          seatElement.innerText = seatNum.toString() + '\n'
          // document.getElementById(elementId).style.backgroundColor = emptySeatColor

          seatElement.onclick = function () {}
        }
      }
    }),
  },
}
</script>

<style scoped></style>
