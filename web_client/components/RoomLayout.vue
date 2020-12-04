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
    // const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_layout')
    // url.search = new URLSearchParams({
    //   room_id: 'english-room',
    // }).toString()
    // const response = await fetch(url.toString(), { method: 'GET' })
    // const res_html = await response.body.room_layout_data
    // console.log('body: ', response.body)
    const res = await common.httpGet(
      'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_layout',
      { room_id: 'english-room' }
    )
    const res_html = res.room_layout_data
    console.log(res_html)
    const dom_parser = new DOMParser()
    try {
      const layout_doc = dom_parser.parseFromString(res_html, 'text/html')
      const roomLayoutDiv = document.getElementById('room-layout')
      await roomLayoutDiv.appendChild(layout_doc.documentElement)
      console.log('どう？')
      console.log(layout_doc.getElementsByClassName('group').length)
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

<style scoped>
#room-layout {
  height: content-box;
  width: content-box;
}
</style>
