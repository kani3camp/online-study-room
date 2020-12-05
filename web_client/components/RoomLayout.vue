<template>
  <div id="room-layout">
    <!--    <Seat-->
    <!--      v-for="seat in layout.seats"-->
    <!--      :key="seat.id"-->
    <!--      :x="seat.x"-->
    <!--      :y="seat.y"-->
    <!--    />-->
    <div
      v-for="(seat, index) in layout.seats"
      :key="seat.id"
      class="seat"
      :style="{ backgroundColor: emptySeatColor, top: seatPositions[index].x+'%', left: seatPositions[index].y+'%', width: seatShape.width+'%' }"
    >
      {{ seat.id }}
    </div>
  </div>
</template>

<script>
import common from '@/plugins/common'
import Seat from '@/components/Seat'

export default {
  name: 'RoomLayout',
  comments: {
    Seat,
  },
  props: [],
  data: () => ({
    emptySeatColor: '#c36b2b',
    layout: {
      room_id: 'mathematics-room',
      'layout-version': 1,
      'num-seats': 3,
      'room-shape': {
        width: 1000,
        height: 800,
      },
      'seat-shape': {
        width: 120,
        height: 100,
      },
      seats: [
        {
          id: 1,
          x: 0,
          y: 100,
        },
        {
          id: 2,
          x: 0,
          y: 200,
        },
        {
          id: 3,
          x: 500,
          y: 0,
        },
      ],
      partitions: [
        {
          id: 1,
          x: 120,
          y: 80,
        },
        {
          id: 2,
          x: 0,
          y: 200,
        },
        {
          id: 3,
          x: 120,
          y: 280,
        },
        {
          id: 4,
          x: 0,
          y: 400,
        },
        {
          id: 5,
          x: 120,
          y: 480,
        },
        {
          id: 6,
          x: 0,
          y: 600,
        },
      ],
    },
  }),
  computed: {
    seatShape: function () {
      return {
        width: (100 * this.layout['seat-shape'].width) / this.layout['room-shape'].width,
        height: (100 * this.layout['seat-shape'].height) / this.layout['room-shape'].height,
      }
    },
    seatPositions: function () {
      const vm = this
      return this.layout.seats.map(function (seat) {
        return vm.seatPosition(seat.x, seat.y)
      })
    },
  },
  async created() {},
  methods: {
    seatPosition(x, y) {
      return {
        x: (100 * x) / this.layout['room-shape'].width,
        y: (100 * y) / this.layout['room-shape'].height,
      }
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
  /*todo*/
  width: 1000px;
  height: 800px;
  background-color: azure;
  margin: auto;
}

.seat {
  width: 4rem;
  height: 3rem;
  position: relative;
}
</style>
