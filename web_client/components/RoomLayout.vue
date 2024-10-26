<template>
  <div>
    <v-card class="mt-3">
      <div
        id="room-layout"
        ref="roomLayout"
        :style="{
          width: roomShape.width,
          height: roomShape.height,
        }"
      >
        <div
          v-for="(seat, index) in seats"
          :key="seat.id"
          class="seat"
          :style="{
            backgroundColor: seat.is_vacant ? emptySeatColor : filledSeatColor,
            left: seatPositions[index].x + '%',
            top: seatPositions[index].y + '%',
            width: seatShape.width + '%',
            height: seatShape.height + '%',
            fontSize: seatFontSize,
          }"
        >
          {{ seat.id + '\n' }}
        </div>

        <div
          v-for="(partition, index) in partitions"
          :key="partition.id"
          class="partition"
          :style="{
            left: partitionPositions[index].x + '%',
            top: partitionPositions[index].y + '%',
            width: partitionShapes[index].width + '%',
            height: partitionShapes[index].height + '%',
          }"
        />
      </div>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { defineProps, computed, watchEffect, useTemplateRef } from 'vue'

const refRoomLayout = useTemplateRef('roomLayout')

const isMounted = ref(false)
const emptySeatColor = ref('#fce7d2')
const filledSeatColor = ref('#430308')
const seats_if_filled = ref(null)
const seats = ref([])
const partitions = ref([])

const props = defineProps({
  layout: {
    type: Object,
    default: null,
  },
})

const roomLayout = computed({
  get: () => props.layout,
  set: (value) => {},
})

const seatFontSize = computed({
  get: () => {
    if (isMounted && roomLayout && refRoomLayout) {
      const roomLayoutWidth = refRoomLayout.value!.clientWidth
      return roomLayoutWidth * roomLayout['font_size_ratio'] + 'px'
    } else {
      return 20 + 'px'
    }
  },
  set: () => {},
})

const roomShape = computed({
  get: () => {
    if (isMounted && roomLayout && refRoomLayout) {
      const roomLayoutWidth = refRoomLayout.value!.clientWidth
      return {
        width: 100 + '%',
        height:
          (roomLayoutWidth * roomLayout['room_shape'].height) / roomLayout['room_shape'].width +
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
  set: () => {},
})

const seatShape = computed({
  get: () => {
    if (this.roomLayout) {
      const vm = this
      // 別なところでこれらの値+'%'で使う？
      return {
        width: (100 * roomLayout['seat_shape'].width) / roomLayout['room_shape'].width,
        height: (100 * roomLayout['seat_shape'].height) / roomLayout['room_shape'].height,
      }
    } else {
      return {
        width: 0,
        height: 0,
      }
    }
  },
  set: () => {},
})

const seatPositions = computed({
  get: () => {
    if (roomLayout) {
      return roomLayout.seats.map((seat) => {
        return {
          x: (100 * seat.x) / props.layout['room_shape'].width,
          y: (100 * seat.y) / props.layout['room_shape'].height,
        }
      })
    } else {
      return []
    }
  },
  set: () => {},
})

const partitionShapes = computed({
  get: () => {
    if (roomLayout) {
      const vm = this
      return roomLayout.partitions.map(function (partition, index) {
        const partitionShapes = roomLayout['partition_shapes']
        const shapeType = partition['shape_type']
        let width
        let height
        for (let i = 0; i < partitionShapes.length; i++) {
          if (partitionShapes[i].name === shapeType) {
            width = (100 * partitionShapes[i].width) / roomLayout['room_shape'].width
            height = (100 * partitionShapes[i].height) / roomLayout['room_shape'].height
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
  set: () => {},
})

const partitionPositions = computed(() => {
  if (roomLayout) {
    const vm = this
    return roomLayout.partitions.map(function (partition) {
      return {
        x: (100 * partition.x) / props.layout['room_shape'].width,
        y: (100 * partition.y) / props.layout['room_shape'].height,
      }
    })
  } else {
    return []
  }
})

onMounted(() => {
  isMounted.value = true
})

const initializeLayoutData = () => {
  // template内で使うときにroomLayout.seatsとして使うと最初はroomLayout = nullなのでエラーになる
  // そのため、個別に変数として扱って、初期値として[]をいれておく
  seats.value = roomLayout.seats
  partitions.value = roomLayout.partitions
}

watchEffect(() => {
  if (Object.keys(roomLayout.value).length !== 0) {
    initializeLayoutData()
  }
})
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
