<template>
  <div id="time">
    <p>{{ year }}年 {{ month }}月 {{ day }}日</p>
    <p class="clock">{{ hours }}：{{ minutes }}</p>
  </div>
</template>

<script>
export default {
  name: 'Time',
  data: () => ({
    date: new Date()
  }),
  computed: {
    year () {
      return this.date.getFullYear()
    },
    month () {
      return this.date.getMonth() + 1
    },
    day () {
      return this.date.getDate()
    },
    hours () {
      return this.date.getHours()
    },
    minutes () {
      let min = this.date.getMinutes().toString()
      if (min.length === 1) {
        min = '0' + min
      }
      return min
    }
  },
  created () {
    setInterval(this.updateTime, 1000)
  },
  methods: {
    updateTime () {
      this.date = new Date()
    }
  }
}
</script>

<style scoped>
#time {
  padding-top: 1rem;
  font-size: larger;
  background-color: aquamarine;
  width: 384px;
}

.clock {
  font-size: 1.5rem;
  font-weight: bold;
}
</style>
