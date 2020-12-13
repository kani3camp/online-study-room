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
            backgroundColor: seat.is_vacant ? emptySeatColor : filledSeatColor,
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
export default {
  name: 'RoomLayout',
  props: {
    layout: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      emptySeatColor: '#fce7d2',
      filledSeatColor: '#430308',
      isMounted: false,
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
    seatFontSize: {
      get() {
        if (this.isMounted && this.roomLayout && this.$refs.roomLayout) {
          const roomLayoutWidth = this.$refs.roomLayout.clientWidth
          return roomLayoutWidth * this.roomLayout['font_size_ratio'] + 'px'
        } else {
          return 20 + 'px'
        }
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
          }
        } else {
          // そもそもlayoutが読み込めてないときは親のコンポーネントで"Loading..."とか表示しておくのでどうでもいい
          return {
            width: 100 + '%',
            height: 100 + '%',
          }
        }
      },
      set() {},
    },
    seatShape: {
      get() {
        if (this.roomLayout) {
          const vm = this
          // 別なところでこれらの値+'%'で使う？
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
            return {
              x: (100 * seat.x) / vm.layout['room_shape'].width,
              y: (100 * seat.y) / vm.layout['room_shape'].height,
            }
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
          return vm.roomLayout.partitions.map(function (partition, index) {
            const partitionShapes = vm.roomLayout['partition_shapes']
            const shapeType = partition['shape_type']
            let width
            let height
            for (let i = 0; i < partitionShapes.length; i++) {
              if (partitionShapes[i].name === shapeType) {
                width = (100 * partitionShapes[i].width) / vm.roomLayout['room_shape'].width
                height = (100 * partitionShapes[i].height) / vm.roomLayout['room_shape'].height
              }
            }
            return {
              width,
              height,
            }
          })
        } else {
          return []
        }
      },
      set() {},
    },
    partitionPositions: function () {
      if (this.roomLayout) {
        const vm = this
        return this.roomLayout.partitions.map(function (partition) {
          return {
            x: (100 * partition.x) / vm.layout['room_shape'].width,
            y: (100 * partition.y) / vm.layout['room_shape'].height,
          }
        })
      } else {
        return []
      }
    },
  },
  watch: {
    roomLayout: function (newValue, oldValue) {
      if (newValue !== {}) {
        this.initializeLayoutData()
      }
    },
  },
  mounted() {
    this.isMounted = true
  },
  methods: {
    initializeLayoutData() {
      // template内で使うときにroomLayout.seatsとして使うと最初はroomLayout = nullなのでエラーになる
      // そのため、個別に変数として扱って、初期値として[]をいれておく
      this.seats = this.roomLayout.seats
      this.partitions = this.roomLayout.partitions
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
  max-width: 800px;
  max-height: 90vh;
}
#room-layout:before {
  content: '';
  display: block;
}

.seat {
  position: absolute;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.partition {
  position: absolute;
  background-color: #647a7f;
}
</style>










