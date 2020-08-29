<template>
  <div id="time">
    time
    {{ date }}
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
      return this.date().getFullYear()
    },
    month () {
      return this.date().getMonth()
    },
    day () {
      return this.date().getDay()
    },
    hours () {
      return this.date().getHours()
    },
    minutes () {
      return this.date().getMinutes()
    },
    seconds () {
      return this.date()
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
  background-color: aquamarine;
  width: 30vw;
}
</style>
