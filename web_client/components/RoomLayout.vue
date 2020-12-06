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
        fontSize: seatFontSize
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
import roomLayoutJson from 'assets/mathematics-rom-layout.json'

export default {
  name: 'RoomLayout',
  comments: {
    Seat,
  },
  props: [],
  data() {
    return {
      emptySeatColor: '#e5b796',
      // seatFontSize: 80 + '%',
      seatFontSize: 80 + '%',
      // layout: {
      //   room_id: 'mathematics-room',
      //   version: 1,
      //   num_seats: 3,
      //   room_shape: {
      //     width: 1000,
      //     height: 800,
      //   },
      //   seat_shape: {
      //     width: 120,
      //     height: 100,
      //   },
      //   partition_shapes: [
      //     {
      //       name: 'vertical=seat',
      //       width: 10,
      //       height: 100,
      //     },
      //     {
      //       name: 'horizontal=seat',
      //       width: 120,
      //       height: 10,
      //     },
      //     {
      //       name: 'vertical=seat+margin',
      //       width: 10,
      //       height: 120,
      //     },
      //     {
      //       name: 'horizontal=seat+margin',
      //       width: 140,
      //       height: 10,
      //     },
      //     {
      //       name: 'vertical=seat+partition-height',
      //       width: 10,
      //       height: 110,
      //     },
      //     {
      //       name: 'horizontal=seat+partition-width',
      //       width: 130,
      //       height: 10,
      //     },
      //   ],
      //   seats: [
      //     { id: 1, x: 0, y: 100 },
      //     { id: 2, x: 0, y: 300 },
      //     { id: 3, x: 0, y: 500 },
      //     { id: 4, x: 0, y: 700 },
      //   ],
      //   partitions: [
      //     { id: 101, x: 120, y: 80, shape_type: 'vertical=seat+margin' },
      //     { id: 102, x: 0, y: 200, shape_type: 'horizontal=seat+partition-width' },
      //     { id: 103, x: 120, y: 280, shape_type: 'vertical=seat+margin' },
      //     { id: 104, x: 0, y: 400, shape_type: 'horizontal=seat+partition-width' },
      //     { id: 105, x: 120, y: 480, shape_type: 'vertical=seat+margin' },
      //     { id: 106, x: 0, y: 600, shape_type: 'horizontal=seat+partition-width' },
      //     { id: 107, x: 120, y: 680, shape_type: 'vertical=seat+margin' },
      //   ],
      // },
      // layout: null,
      layout: roomLayoutJson,
      isMounted: false,
    }
  },
  computed: {
    // seatFontSize: {
    //   get() {
    //     return this.roomShape.width + 'px'
    //   },
    //   set(value) {},
    // },
    roomShape: {
      get() {
        if (this.isMounted) {
          console.log('再度')
          const roomLayoutWidth = this.$refs.roomLayout.clientWidth
          const roomLayoutHeight = this.$refs.roomLayout.clientHeight
          if (roomLayoutWidth > roomLayoutHeight) {
            // 横長画面
            const widthPx =
              (roomLayoutHeight * this.layout['room_shape'].width) / this.layout['room_shape'].height
            return {
              width: widthPx + 'px',
              height: roomLayoutHeight + 'px',
              seatFontSize: widthPx * this.layout['font_size_ratio'] + 'px',
            }
          } else {
            // 縦長画面
            return {
              width: roomLayoutWidth + 'px',
              height:
                (roomLayoutWidth * this.layout['room_shape'].height) / this.layout['room_shape'].width + 'px',
            }
          }
        } else {
          return {
            width: 100 + 'vw',
            height: 90 + 'vh',
          }
        }
      },
      set() {},
    },
    seatShape: {
      get() {
        const vm = this
        return {
          width: (100 * vm.layout['seat_shape'].width) / vm.layout['room_shape'].width,
          height: (100 * vm.layout['seat_shape'].height) / vm.layout['room_shape'].height,
        }
      },
      set() {},
    },
    seatPositions: {
      get() {
        const vm = this
        return this.layout.seats.map(function (seat) {
          return vm.seatPosition(seat.x, seat.y)
        })
      },
      set() {},
    },
    partitionShapes: {
      get() {
        const vm = this
        return vm.layout.partitions.map(function (partition, index) {
          const shapeType = vm.layout.partitions[index]['shape_type']
          let width
          let height
          for (let i = 0; i < vm.layout['partition_shapes'].length; i++) {
            if (vm.layout['partition_shapes'][i].name === shapeType) {
              width = (100 * vm.layout['partition_shapes'][i].width) / vm.layout['room_shape'].width
              height = (100 * vm.layout['partition_shapes'][i].height) / vm.layout['room_shape'].height
            }
          }
          return {
            width: width,
            height: height,
          }
        })
      },
      set() {},
    },
    partitionPositions: function () {
      const vm = this
      return this.layout.partitions.map(function (partition) {
        return vm.seatPosition(partition.x, partition.y)
      })
    },
  },
  watch: {
    roomShape: function (newValue, oldValue) {
      this.seatFontSize = newValue.seatFontSize
    },
  },
  created() {
    console.log(this.layout)
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
      this.seatFontSize = Math.floor((roomLayoutWidth * 30) / 800)
      this.seatShape = {
        width: Math.floor(
          (roomLayoutWidth * this.layout['seat_shape'].width) / this.layout['room_shape'].width
        ),
        height: Math.floor(
          (roomLayoutHeight * this.layout['seat_shape'].height) / this.layout['room_shape'].height
        ),
      }
      this.seatPositions = this.layout.seats.map(function (seat) {
        return {
          x: (roomLayoutWidth * seat.x) / vm.layout['room_shape'].width,
          y: (roomLayoutWidth * seat.y) / vm.layout['room_shape'].height,
        }
      })
    },
    seatPosition(x, y) {
      return {
        x: (100 * x) / this.layout['room_shape'].width,
        y: (100 * y) / this.layout['room_shape'].height,
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
  background-color: azure;
  margin: auto;
  border: solid 0.5rem black;
  max-width: 100vw;
  max-height: 90vh;
}
#room-layout:before {
  content: '';
  display: block;
  /*padding-top: todo ;*/
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










