<template>
  <div>
    <v-card
      class="mt-3"
    >

      <div
        id="room-layout"
        ref="roomLayout"
        :style="{
          width: roomShape.width,
          height: roomShape.height
        }"
      >
        <div
          v-for="(seat, index) in seats"
          :key="seat.id"
          class="seat"
          :style="{
            backgroundColor: seats_if_filled[seat.id] ? filledSeatColor : emptySeatColor,
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
          v-for="(partition, index) in partitions"
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
    </v-card>
  </div>
</template>

<script>
import Seat from '@/components/Seat'
// todo コード整理
export default {
  name: 'RoomLayout',
  comments: {
    Seat,
  },
  props: {
    roomId: {
      type: String,
      required: true,
    },
    layout: {
      type: Object,
      default: null,
    },
    seatsDataArray: {
      type: Array,
      default: null,
    },
  },
  data() {
    return {
      emptySeatColor: '#eaccb6',
      filledSeatColor: '#431e03',
      seatFontSize: 80 + '%',
      isMounted: false,
      // isLayoutLoaded: false,
      seats_if_filled: null,
      seats: [],
      partitions: [],
    }
  },
  computed: {
    roomLayout: {
      get() {
        return this.layout
      },
      set() {},
    },
    seatsData: {
      get() {
        return this.seatsDataArray
      },
      set() {},
    },
    roomShape: {
      get() {
        if (this.isMounted && this.roomLayout && this.$refs.roomLayout) {
          const roomLayoutWidth = this.$refs.roomLayout.clientWidth
          return {
            width: 100 + '%',
            height:
              (roomLayoutWidth * this.roomLayout['room_shape'].height) / this.roomLayout['room_shape'].width +
              'px',
            seatFontSize: roomLayoutWidth * this.roomLayout['font_size_ratio'] + 'px',
          }
        } else {
          return {
            width: 100 + '%',
            height: 100 + 'vh',
          }
        }
      },
      set() {},
    },
    seatShape: {
      get() {
        if (this.roomLayout) {
          const vm = this
          return {
            width: (100 * vm.roomLayout['seat_shape'].width) / vm.roomLayout['room_shape'].width,
            height: (100 * vm.roomLayout['seat_shape'].height) / vm.roomLayout['room_shape'].height,
          }
        } else {
          return {
            width: 0,
            height: 0,
          }
        }
      },
      set() {},
    },
    seatPositions: {
      get() {
        if (this.roomLayout) {
          const vm = this
          return this.roomLayout.seats.map(function (seat) {
            return vm.seatPosition(seat.x, seat.y)
          })
        } else {
          return []
        }
      },
      set() {},
    },
    partitionShapes: {
      get() {
        if (this.roomLayout) {
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
        } else {
          return {
            width: 0,
            height: 0,
          }
        }
      },
      set() {},
    },
    partitionPositions: function () {
      if (this.roomLayout) {
        const vm = this
        return this.layout.partitions.map(function (partition) {
          return vm.seatPosition(partition.x, partition.y)
        })
      } else {
        return []
      }
    },
    roomSize: function () {
      if (this.$refs.roomLayout) {
        const roomLayoutWidth = this.$refs.roomLayout.clientWidth
        const roomLayoutHeight = this.$refs.roomLayout.clientHeight
        return {
          width: roomLayoutWidth,
          height: roomLayoutHeight,
        }
      }
      return {
        width: 0,
        height: 0,
      }
    },
  },
  watch: {
    roomLayout: function (newValue, oldValue) {
      if (newValue !== {}) {
        this.initializeLayoutData()
        this.determineFontSize()
      }
    },
    seatsData: function (newValue, oldValue) {
      if (newValue && this.seats) {
        for (let seat in this.seats) {
          // this.seats_if_filled[seat.id] = false
        }
        for (let filled_seat in newValue) {
          // this.seats_if_filled[filled_seat.seat_id] = true
        }
      }
    },
    roomShape: function (newValue, oldValue) {
      if (this.roomLayout) {
        this.seatFontSize = newValue.seatFontSize
      }
    },
    roomSize: function (newValue, oldValue) {
      if (newValue.width !== 0) {
      }
    },
  },
  async mounted() {
    this.isMounted = true
  },
  methods: {
    initializeLayoutData() {
      this.seats = this.roomLayout.seats
      for (let seat of this.seats) {
        console.log(seat)
        this.seats_if_filled[seat.id] = false
      }
      this.partitions = this.roomLayout.partitions
    },
    determineFontSize() {
      const roomLayoutWidth = this.$refs.roomLayout.clientWidth
      const roomLayoutHeight = this.$refs.roomLayout.clientHeight
      const vm = this
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










