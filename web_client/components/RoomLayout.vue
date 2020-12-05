<template>
  <div
    id="room-layout"
    ref="roomLayout"
    :style="{
      width: roomShape.width,
      height: roomShape.height
    }"
  >
    <div
      v-for="(seat, index) in layout.seats"
      :key="seat.id"
      class="seat"
      :style="{
        backgroundColor: emptySeatColor,
        left: seatPositions[index].x+'%',
        top: seatPositions[index].y+'%',
        width: seatShape.width+'%',
        height: seatShape.height+'%',
        fontSize: defaultSeatFontSize+'px'
      }"
    >
      {{ seat.id + '\n' }}
    </div>

    <div
      v-for="(partition, index) in layout.partitions"
      :key="partition.id"
      class="partition"
      :style="{
        left: partitionPositions[index].x+'%',
        top: partitionPositions[index].y+'%',
        width: partitionShapes[index].width+'%',
        height: partitionShapes[index].height+'%'
      }"
    />
  </div>
</template>

<script>
import Seat from '@/components/Seat'

export default {
  name: 'RoomLayout',
  comments: {
    Seat,
  },
  props: [],
  data() {
    return {
      emptySeatColor: '#e5b796',
      defaultSeatFontSize: 0,
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
        'partition-shapes': {
          'vertical=seat': {
            width: 10,
            height: 100,
          },
          'horizontal=seat': {
            width: 120,
            height: 10,
          },
          'vertical=seat+margin': {
            width: 10,
            height: 120,
          },
          'horizontal=seat+margin': {
            width: 140,
            height: 10,
          },
          'vertical=seat+partition-height': {
            width: 10,
            height: 110,
          },
          'horizontal=seat+partition-width': {
            width: 130,
            height: 10,
          },
        },
        seats: [
          { id: 1, x: 0, y: 100 },
          { id: 2, x: 0, y: 300 },
          { id: 3, x: 0, y: 500 },
          { id: 4, x: 0, y: 700 },
        ],
        partitions: [
          { id: 101, x: 120, y: 80, shapeType: 'vertical=seat+margin' },
          { id: 102, x: 0, y: 200, shapeType: 'horizontal=seat+partition-width' },
          { id: 103, x: 120, y: 280, shapeType: 'vertical=seat+margin' },
          { id: 104, x: 0, y: 400, shapeType: 'horizontal=seat+partition-width' },
          { id: 105, x: 120, y: 480, shapeType: 'vertical=seat+margin' },
          { id: 106, x: 0, y: 600, shapeType: 'horizontal=seat+partition-width' },
          { id: 107, x: 120, y: 680, shapeType: 'vertical=seat+margin' },
        ],
      },
      isMounted: false,
    }
  },
  computed: {
    roomShape: function () {
      if (this.isMounted) {
        console.log('再度')
        const roomLayoutWidth = this.$refs.roomLayout.clientWidth
        return {
          width: 100 + '%',
          height:
            (roomLayoutWidth * this.layout['room-shape'].height) / this.layout['room-shape'].width + 'px',
        }
      } else {
        return {
          width: 100 + '%',
          height: this.layout['room-shape'].width + 'px',
        }
      }
    },
    seatShape: function () {
      const vm = this
      return {
        width: (100 * vm.layout['seat-shape'].width) / vm.layout['room-shape'].width,
        height: (100 * vm.layout['seat-shape'].height) / vm.layout['room-shape'].height,
      }
      // if (this.isMounted) {
      //   console.log('mountされてるね')
      //   const roomLayoutWidth = this.$refs.roomLayout.clientWidth
      //   const roomLayoutHeight = this.$refs.roomLayout.clientHeight
      //   return {
      //     width: (roomLayoutWidth * vm.layout['seat-shape'].width) / vm.layout['room-shape'].width,
      //     height: (roomLayoutHeight * vm.layout['seat-shape'].height) / vm.layout['room-shape'].height,
      //   }
      // } else {
      //   return {
      //     width: (100 * this.layout['seat-shape'].width) / this.layout['room-shape'].width,
      //     height: (100 * this.layout['seat-shape'].height) / this.layout['room-shape'].height,
      //   }
      // }
    },
    seatPositions: function () {
      const vm = this
      return this.layout.seats.map(function (seat) {
        return vm.seatPosition(seat.x, seat.y)
      })
      // if (vm.isMounted) {
      // const roomLayoutWidth = this.$refs.roomLayout.clientWidth
      // const roomLayoutHeight = this.$refs.roomLayout.clientHeight
      // return vm.layout.seats.map(function (seat) {
      //   return {
      //     x: (roomLayoutWidth * seat.x) / vm.layout['room-shape'].width,
      //     y: (roomLayoutHeight * seat.y) / vm.layout['room-shape'].height,
      //   }
      // })
      // } else {
      //   return this.layout.seats.map(function (seat) {
      //     return vm.seatPosition(seat.x, seat.y)
      //   })
      // }
    },
    partitionShapes: function () {
      const vm = this
      return vm.layout.partitions.map(function (partition, index) {
        const shapeType = vm.layout.partitions[index]['shapeType']
        return {
          width: (100 * vm.layout['partition-shapes'][shapeType].width) / vm.layout['room-shape'].width,
          height: (100 * vm.layout['partition-shapes'][shapeType].height) / vm.layout['room-shape'].height,
        }
      })
    },
    partitionPositions: function () {
      const vm = this
      return this.layout.partitions.map(function (partition) {
        return vm.seatPosition(partition.x, partition.y)
      })
    },
  },
  async mounted() {
    console.log('mounted()')
    this.isMounted = true
    this.determineFontSize()
  },
  methods: {
    determineFontSize() {
      const roomLayoutWidth = this.$refs.roomLayout.clientWidth
      const roomLayoutHeight = this.$refs.roomLayout.clientHeight
      const vm = this
      console.log(roomLayoutWidth)
      this.defaultSeatFontSize = Math.floor((roomLayoutWidth * 30) / 800)
      this.seatShape = {
        width: Math.floor(
          (roomLayoutWidth * this.layout['seat-shape'].width) / this.layout['room-shape'].width
        ),
        height: Math.floor(
          (roomLayoutHeight * this.layout['seat-shape'].height) / this.layout['room-shape'].height
        ),
      }
      this.seatPositions = this.layout.seats.map(function (seat) {
        return {
          x: (roomLayoutWidth * seat.x) / vm.layout['room-shape'].width,
          y: (roomLayoutWidth * seat.y) / vm.layout['room-shape'].height,
        }
      })
    },
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
  position: relative;
  top: 0;
  left: 0;
  box-sizing: border-box;
  background-color: azure;
  margin: auto;
}

.seat {
  position: absolute;
  display: flex;
  justify-content: center;
  align-items: center;
}

.partition {
  position: absolute;
  background-color: #647a7f;
}
</style>
