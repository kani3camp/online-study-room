<template>
  <div
    :id="seatId"
    class="seat"
    :style="{ backgroundColor: emptySeatColor }"
  />
</template>

<script>
export default {
  name: 'Seat',
  props: {
    x: {
      type: Number,
      default: 0,
    },
    y: {
      type: Number,
      default: 0,
    },
    id: {
      type: Number,
      default: 0,
    },
  },
  data: () => ({
    emptySeatColor: '#c36b2b',
    maxDisplayNameLength: 6,
    seatId: 'seat-' + this.id.toString(),
  }),
}
</script>

<style scoped>
.seat {
  /*todo*/
  width: 3rem;
  height: 2rem;
  position: relative;
  top: 0;
  left: 0;
}
</style>
