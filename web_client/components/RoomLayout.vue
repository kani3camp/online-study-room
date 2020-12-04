<template>
  <!--  todo このdivは消す-->
  <iframe src="../assets/english-room.html" />
</template>

<script>
export default {
  name: 'RoomLayout',
  props: [],
  data: () => ({
    layoutFile: '@/assets/english-room.html',
    emptySeatColor: '#ddcec3',
    maxDisplayNameLength: 6,
  }),
  mounted() {},
  methods: {
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
