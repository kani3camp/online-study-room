<template :src="layoutFile">
  <!--  todo このdivは消す-->
  <div>
    <input
      type="file"
      @change="onChangeFile"
    >
    <div v-html="layoutRawHtml" />
  </div>
</template>

<script>
export default {
  name: 'RoomLayout',
  props: [],
  data: () => ({
    layoutFile: '@/assets/english-room.html',
    layoutRawHtml: '',
    emptySeatColor: '#ddcec3',
    maxDisplayNameLength: 6,
  }),
  async mounted() {},
  methods: {
    onChangeFile: function (event) {
      const file = event.target.files[0]
      if (!file) {
        console.log('no file')
        return false
      }
      const vm = this
      const reader = new FileReader()
      reader.onload = function (e) {
        vm.layoutRawHtml = e.target.result
        console.log('layoutRawHtml: ', this.layoutRawHtml)
      }
      reader.readAsText(file)
    },
    window: (onload = function () {
      const numSeats = 23

      for (let n = 0; n < numSeats; n++) {
        const seatNum = n + 1
        const seatElement = document.getElementById('seat-' + seatNum.toString())
        seatElement.innerText = seatNum.toString() + '\n'
        // document.getElementById(elementId).style.backgroundColor = emptySeatColor

        seatElement.onclick = function () {}
      }
    }),
  },
}
</script>

<style scoped></style>
