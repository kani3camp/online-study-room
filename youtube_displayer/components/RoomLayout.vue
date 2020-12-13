<template>
  <div>
    <div
      id="room-layout"
      ref="roomLayout"
      :style="{
        width: roomShape.widthPx+'px',
        height: roomShape.heightPx+'px'
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
          fontSize: seatFontSizePx+'px'
        }"
      >
        <p style="font-weight: bold">
          {{ seat.id }}
        </p>
        <p style="overflow: hidden;">
          {{ seat.user_name.substr(0, 4) }}
        </p>
      </div>

      <div
        v-for="(partition, index) in partitions"
        :key="partition.id"
        class="partition"
        :style="{
          left: partitionPositions[index].x+'%',
          top: partitionPositions[index].y+'%',
          width: partitionShapes[index].widthPercent+'%',
          height: partitionShapes[index].heightPercent+'%'
        }"
      />
    </div>
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
  data () {
    return {
      emptySeatColor: '#fce7d2',
      filledSeatColor: '#ee989f',
      isMounted: false,
      seats_if_filled: null,
    }
  },
  computed: {
    roomLayout: {
      get () {
        return this.layout
      },
      set () {},
    },
    userNamesMap: {
      get () {
        if (this.userNamesMapProp) {
          return this.userNamesMapProp
        } else {
          return null
        }
      },
      set () {},
    },
    seatFontSizePx: {
      get () {
        if (this.isMounted && this.roomLayout && this.roomShape) {
          const roomLayoutWidth = this.roomShape.width
          return roomLayoutWidth * this.roomLayout.font_size_ratio
        } else {
          return 30
        }
      },
      set () {},
    },
    roomShape: {
      get () {
        if (this.isMounted && this.roomLayout) {
          return {
            widthPx: 800 * this.roomLayout.room_shape.width / this.roomLayout.room_shape.height,
            heightPx: 800,
          }
        } else {
          // そもそもlayoutが読み込めてないときは親のコンポーネントで"Loading..."とか表示しておくのでどうでもいい
          return {
            widthPx: 1100,
            heightPx: 800,
          }
        }
      },
      set () {},
    },
    seatShape: {
      get () {
        if (this.roomLayout) {
          const vm = this
          // 別なところでこれらの値+'%'で使う？
          return {
            width: (100 * vm.roomLayout.seat_shape.width) / vm.roomLayout.room_shape.width,
            height: (100 * vm.roomLayout.seat_shape.height) / vm.roomLayout.room_shape.height,
          }
        } else {
          return {
            width: 0,
            height: 0,
          }
        }
      },
      set () {},
    },
    seatPositions: {
      get () {
        if (this.roomLayout) {
          const vm = this
          return this.roomLayout.seats.map(seat => ({
            x: (100 * seat.x) / vm.layout.room_shape.width,
            y: (100 * seat.y) / vm.layout.room_shape.height,
          }))
        } else {
          return []
        }
      },
      set () {},
    },
    partitionShapes: {
      get () {
        if (this.roomLayout) {
          const vm = this
          return vm.roomLayout.partitions.map((partition, index) => {
            const partitionShapes = vm.roomLayout.partition_shapes
            const shapeType = partition.shape_type
            let widthPercent
            let heightPercent
            for (let i = 0; i < partitionShapes.length; i++) {
              if (partitionShapes[i].name === shapeType) {
                widthPercent = (100 * partitionShapes[i].width) / vm.roomLayout.room_shape.width
                heightPercent = (100 * partitionShapes[i].height) / vm.roomLayout.room_shape.height
              }
            }
            return {
              widthPercent,
              heightPercent,
            }
          })
        } else {
          return []
        }
      },
      set () {},
    },
    partitionPositions () {
      if (this.roomLayout) {
        const vm = this
        return this.roomLayout.partitions.map(partition => ({
          x: (100 * partition.x) / vm.layout.room_shape.width,
          y: (100 * partition.y) / vm.layout.room_shape.height,
        }))
      } else {
        return []
      }
    },
    seats () {
      if (this.roomLayout) {
        return this.roomLayout.seats
      } else {
        return []
      }
    },
    partitions () {
      if (this.roomLayout) {
        return this.roomLayout.partitions
      } else {
        console.log('もどした')
        return []
      }
    }
  },
  watch: {},
  mounted () {
    this.isMounted = true
  },
  methods: {},
}
</script>

<style scoped>
#room-layout {
  position: relative;
  top: 0;
  left: 0;
  background-color: azure;
  margin: auto;
  border: solid 6px #666666;
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
