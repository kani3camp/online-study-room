<template>
  <div>

    <div id="seat-selector">
      <v-form class="mx-auto">
        <v-select
          v-model="selected_seat_id"
          :items="layout.seats"
          item-value="id"
          item-text="id"
          label="座席番号"
          outlined
        />
        <div>
          <v-btn
            color="primary"
            block
            elevation="3"
            :disabled="! selected_seat_id"
            @click="select"
          >
            決定
          </v-btn>
        </div>
      </v-form>
    </div>

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
    </v-card>
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
  props: {
    roomId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      emptySeatColor: '#e5b796',
      seatFontSize: 80 + '%',
      layout: roomLayoutJson,
      isMounted: false,
      selected_seat_id: null,
    }
  },
  computed: {
    roomShape: {
      // get() {
      //   if (this.isMounted) {
      //     console.log('再度')
      //     const roomLayoutWidth = this.$refs.roomLayout.clientWidth
      //     const roomLayoutHeight = this.$refs.roomLayout.clientHeight
      //     if (roomLayoutWidth > roomLayoutHeight) {
      //       // 横長画面
      //       const widthPx =
      //         (roomLayoutHeight * this.layout['room_shape'].width) / this.layout['room_shape'].height
      //       return {
      //         width: widthPx + 'px',
      //         height: roomLayoutHeight + 'px',
      //         seatFontSize: widthPx * this.layout['font_size_ratio'] + 'px',
      //       }
      //     } else {
      //       // 縦長画面
      //       return {
      //         width: roomLayoutWidth + 'px',
      //         height:
      //           (roomLayoutWidth * this.layout['room_shape'].height) / this.layout['room_shape'].width + 'px',
      //       }
      //     }
      //   } else {
      //     return {
      //       width: 100 + 'vw',
      //       height: 90 + 'vh',
      //     }
      //   }
      // },
      get() {
        if (this.isMounted) {
          console.log('再度roomShape')
          const roomLayoutWidth = this.$refs.roomLayout.clientWidth
          return {
            width: 100 + '%',
            height:
              (roomLayoutWidth * this.layout['room_shape'].height) / this.layout['room_shape'].width + 'px',
            seatFontSize: roomLayoutWidth * this.layout['font_size_ratio'] + 'px',
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
    created() {
      console.log(this.layout)
    },
    async mounted() {
      console.log('mounted()')
      this.isMounted = true
      this.determineFontSize()
    },
  },
  mounted() {
    this.isMounted = true
  },
  methods: {
    select() {
      this.$emit('selected', this.selected_seat_id)
    },
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










